const fs = require('fs');
const yaml = require('yaml');

// This name must match the alias we created in the #oac-automation-reports Slack channel.
const heroAliasName = '@oac-automation-watchers';

let tutoneConfig = null;
let schemaOld = null;
let schemaLatest = null;
let heroMention = "";

const pathPrefix = '';
try {
  const tutoneConfigFile = fs.readFileSync(`${pathPrefix}.tutone.yml`, 'utf8')
  tutoneConfig = yaml.parse(tutoneConfigFile)

  const schemaFileOld = fs.readFileSync(`${pathPrefix}schema-test-old.json`, 'utf8');
  // const schemaFileOld = fs.readFileSync('schema-old.json', 'utf8');
  schemaOld = JSON.parse(schemaFileOld);

  const schemaFileLatest = fs.readFileSync(`${pathPrefix}schema-test-new.json`, 'utf8');
  // const schemaFileLatest = fs.readFileSync('schema.json', 'utf8');
  schemaLatest = JSON.parse(schemaFileLatest);
} catch (err) {
  console.error(err);
}

// Check for any newly added mutations
const endpointsOld = schemaOld.mutationType.fields.map(field => field.name);
const endpointsLatest = schemaLatest.mutationType.fields.map(field => field.name);
const newEndpoints = endpointsLatest.filter(x => !endpointsOld.includes(x));
const hasNewEndpoints = newEndpoints.length > 0;

// console.log('newEndpoints:', newEndpoints);

// Get the mutations the client has implemented
const clientMutations = tutoneConfig.packages.map(pkg => {
  if (!pkg.mutations) {
    return null;
  }

  if (!pkg.mutations.length) {
    return null;
  }

  return pkg.mutations.map(m => m.name)
}).flat().reduce((acc, i) => i ? [...acc, i] : acc, []);

const clientEndpointsSchemaOld = schemaOld.mutationType.fields.filter(field => clientMutations.includes(field.name));
const clientEndpointsSchemaNew = schemaLatest.mutationType.fields.filter(field => clientMutations.includes(field.name));

// Check for changes in the mutations' signatures
const changedEndpoints = clientEndpointsSchemaNew.reduce((arr, field) => {
  const oldMatch = clientEndpointsSchemaOld.find(f => f.name === field.name);
  if (!oldMatch) {
    return [...arr];
  }

  if (!oldMatch.args?.length && !field.args?.length) {
    return [...arr];
  }

  const differences = compareArrays(oldMatch.args, field.args);
  if (differences.length) {
    return [...arr, {
      name: field.name,
      diff: differences,
    }];
  }

  return [...arr];
}, []);

// Generates a package name based on the endpoint name.
// If an endpoint contains a substring of the keywords listed below,
// it takes the substring leading up to that to generate the package name
// (honestly this isn't fully reliable bc endpoint names are not fully
// compliant with standards).
function generatePackageNameForEndpoint(endpointName) {
  const keywords = /Create|Read|Update|Delete|Add|Remove|Revoke|Write/;

  return endpointName.split(keywords)[0].toLowerCase();
}

// console.log('');

const packages = [];

const newMutationsConfigs = newEndpoints.map((endpointName) => {
  const newEndpointSchema = schemaLatest.mutationType.fields.find(f => f.name === endpointName);
  const args = newEndpointSchema.args
  const pkgName = generatePackageNameForEndpoint(endpointName);

  if (!packages.includes(pkgName)) {
    packages.push({
      name: pkgName,
      import_path: `github.com/newrelic/newrelic-client-go/v2/pkg/${pkgName}`,
      generators: ["typegen", "nerdgraphclient"],
      mutations: [],
    });
  }

  let maxQueryDepth = 1;
  for (let i = 0; i < args.length; i++) {
    maxQueryDepth = getMaxQueryDepth(args[i].type);
  }

  // console.log('pkgName:', pkgName);

  return {
    packageName: pkgName,
    name: endpointName,
    maxQueryDepth,
  };
});

// console.log('newMutationsConfigs:', newMutationsConfigs);

function mergeObjectsArray(objects) {
  return objects.reduce((acc, current) => {
    // Check if the object already exists in the accumulator
    if (!acc.some(obj => obj.name === current.name)) {
      acc.push(current);
    }

    return acc;
  }, []);
}

function findPackageByName(packages, packageName) {
  return packages.find(pkg => pkg.name === packageName);
}

// Create the config to generate the feature code.
const config = mergeObjectsArray(newMutationsConfigs.reduce((arr, mutationConfig) => {
  const pkg = findPackageByName(packages, mutationConfig.packageName);

  // If no package, it's new, so make a new package config
  // TODO: Need to make sure we don't overwrite existing packages
  //       in the tutone config. This is currently a naive implementation.
  if (!pkg) {
    return arr;
  }

  pkg.mutations = [...pkg.mutations, {
    name: mutationConfig.name,
    maxQueryDepth: mutationConfig.maxQueryDepth,
  }];

  return [...arr, pkg];
}, []));

let cfg = {
  log_level: 'trace',
  cache: {
    schema_file: 'schema.json',
  },
  endpoint: 'https://api.newrelic.com/graphql',
  auth: {
    header: "Api-Key",
    api_key_env_var: 'NEW_RELIC_API_KEY',
  },
};

cfg.packages = config;

const tutoneConfigYAML = yaml.stringify(cfg);

// console.log('packages:', config);
// console.log('newMutationsConfig:', newMutationsConfigs);


// console.log('config JSON:\n\n', JSON.stringify(config, null, 2));
console.log('');
console.log('Tutone config:');
console.log('');
console.log(tutoneConfigYAML);
// console.log('');

// Check to see which mutations the client is missing
const schemaMutations = schemaLatest.mutationType.fields.map(field => field.name);
const clientMutationsDiff = schemaMutations.filter(x => !clientMutations.includes(x));

// console.log('');
// console.log('Client is still missing the following API mutations:\n', clientMutationsDiff);
// console.log('');

let newApiMutationsMsg = '';
if (hasNewEndpoints) {
  heroMention = heroAliasName;
  newApiMutationsMsg = `'${newEndpoints.join('\n')}'`;
}


let clientMutationsDiffMsg = ''
if (clientMutationsDiff.length > 0) {
  clientMutationsDiffMsg = `'${clientMutationsDiff.join('\n')}'`;
}

function compareObjects(obj1, obj2, path = '') {
  let differences = [];

  for (let key in obj1) {
    if (!obj2.hasOwnProperty(key)) {
      differences.push({
        old: obj1[key],
        new: undefined,
        property: `${path}.${key}`,
      });

      continue;
    }

    if (obj2.hasOwnProperty(key)) {
      if (typeof obj1[key] === 'object' && typeof obj2[key] === 'object') {
        differences = differences.concat(compareObjects(obj1[key], obj2[key], `${path}.${key}`));

        continue;
      }

      if (obj1[key] !== obj2[key]) {
        differences.push({
          old: obj1[key],
          new: obj2[key],
          property: `${path}.${key}`,
        });

        continue;
      }
    }
  }

  for (let key in obj2) {
    if (obj2.hasOwnProperty(key) && !obj1.hasOwnProperty(key)) {
      differences.push({
        old: undefined,
        new: obj2[key],
        property: `${path}.${key}`,
      });
    }
  }

  return differences;
}

function compareArrays(arr1, arr2) {
  let differences = [];

  for (let i = 0; i < Math.max(arr1.length, arr2.length); i++) {
    if (arr1[i] && arr2[i]) {
      differences = differences.concat(compareObjects(arr1[i], arr2[i], `[${i}]`));
    } else if (arr1[i]) {
      differences.push({
        property: `[${i}]`,
        old: arr1[i],
        new: null,
      });
    } else if (arr2[i]) {
      differences.push({
        property: `[${i}]`,
        old: null,
        new: arr2[i]
      });
    }
  }

  return differences;
}

function getTypeFromSchema(typeName) {
  return schemaLatest.types.find(t => t.name === typeName);
}

function getTypeName(type) {
  if (type.ofType) {
    // Recursion FTW
    return getTypeName(type.ofType);
  }

  if (type.name !== "") {
    return type.name;
  }
}

function getMaxQueryDepth(type, depth = 1) {
  let maxQueryDepth = depth;

  type = getTypeFromSchema(getTypeName(type));

  if (type?.kind === 'INPUT_OBJECT' || type?.ofType?.kind === 'INPUT_OBJECT') {
    maxQueryDepth++

    for (const field of type.inputFields) {
      if (field.type?.kind === 'INPUT_OBJECT' || field.type.ofType?.kind === 'INPUT_OBJECT') {
        const inputType = getTypeFromSchema(getTypeName(field.type));
        if (inputType && inputType.kind === 'INPUT_OBJECT') {
          // Recursion FTW
          maxQueryDepth = Math.max(maxQueryDepth, getMaxQueryDepth(inputType, depth + 1));
        }
      }

      if (field.type.kind === 'LIST' && field.type?.ofType?.ofType.kind === 'INPUT_OBJECT') {
        const inputType = getTypeFromSchema(getTypeName(field.type.ofType.ofType));

        // Recursion FTW
        maxQueryDepth = Math.max(maxQueryDepth, getMaxQueryDepth(inputType, depth + 1));
      }
    }
  }

  // console.log('maxQueryDepth', maxQueryDepth);

  return maxQueryDepth;
}

module.exports = {
  heroMention,
  schemaMutations,
  clientMutations,
  clientMutationsDiff,
  newApiMutationsMsg,
  clientMutationsDiffMsg,
  changedEndpoints,
  tutoneConfig: tutoneConfigYAML,
};
