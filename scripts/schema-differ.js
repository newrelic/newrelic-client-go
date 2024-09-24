const fs = require('fs');
const yaml = require('yaml');
const merge = require('lodash.merge');

// This name must match the alias we created in the #oac-automation-reports Slack channel.
const heroAliasName = '@oac-automation-watchers';

let tutoneConfig = null;
let schemaOld = null;
let schemaLatest = null;
let heroMention = '';

const pathPrefix = './';
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

const packagesToGenerate = [];

const clientPackages = tutoneConfig.packages;


newEndpoints.forEach((endpointName) => {
  const newEndpointSchema = schemaLatest.mutationType.fields.find(f => f.name === endpointName);
  const args = newEndpointSchema.args
  const pkgName = generatePackageNameForEndpoint(endpointName);
  const existingPackage = findPackageByName(clientPackages, pkgName);
  const cachedPackage = findPackageByName(packagesToGenerate, pkgName);

  let maxQueryDepth = 1;
  for (let i = 0; i < args.length; i++) {
    maxQueryDepth = getMaxQueryDepth(args[i].type);
  }

  // If we've already added package to the array, we need to add the
  // mutation to the package if the mutation hasn't been added yet.
  if (cachedPackage) {
    const cachedMutation = cachedPackage.mutations?.length > 0
      ? findMutationByName(cachedPackage.mutations, endpointName)
      : null;

    if (!cachedMutation) {
      cachedPackage.mutations = [...cachedPackage.mutations, {
        name: endpointName,
        max_query_field_depth: maxQueryDepth,
      }];
    }
  }

  // If the tutone config doesn't have the package defined and we haven't added
  // to the cached packages array, we need to add the package to the packagesToGenerate array.
  if (!existingPackage && !cachedPackage) {
    packagesToGenerate.push({
      name: pkgName,
      import_path: `github.com/newrelic/newrelic-client-go/v2/pkg/${pkgName}`,
      generators: ["typegen", "nerdgraphclient"],
      mutations: [{
        name: endpointName,
        max_query_field_depth: maxQueryDepth,
      }],
    });
  }

  // If the tutone config already has the package defined, but we haven't added it
  // to the packagesToGenerate array, we need to add the package to the packagesToGenerate array,
  // and add the mutation to the package if the mutation hasn't been added yet.
  if (existingPackage && !cachedPackage) {
    // Clone the existing package so we don't mutate the original
    const pkg = JSON.parse(JSON.stringify(existingPackage));

    // Add the package to the list of packages to generate
    packagesToGenerate.push(pkg);

    const cachedMutation = pkg.mutations?.length > 0
      ? findMutationByName(pkg.mutations, endpointName)
      : null;

    // console.log('cachedMutation:', endpointName);

    // Ensure we don't add the same mutation twice
    // TODO: This is a naive implementation. We should check if
    // the mutation's query depth has changed.
    if (!cachedMutation) {
      pkg.mutations = [...pkg.mutations, {
        name: endpointName,
        max_query_field_depth: maxQueryDepth,
      }];
    }
  }
});

// console.log('packagesToGenerate:', JSON.stringify(packagesToGenerate, null, 2));

function findPackageByName(packages, packageName) {
  return packages.find(pkg => pkg.name === packageName);
}

function findMutationByName(mutations, mutationName) {
  return mutations.find(m => m.name === mutationName);
}

let cfg = {
  log_level: 'info',
  cache: {
    schema_file: 'schema.json',
  },
  endpoint: 'https://api.newrelic.com/graphql',
  auth: {
    header: "Api-Key",
    api_key_env_var: 'NEW_RELIC_API_KEY',
  },
  generators: [
    {
      name: "typegen",
      fileName: "types.go"
    },
    {
      name: "nerdgraphclient",
      fileName: "{{.PackageName}}_api.go"
    },
  ]
};

cfg.packages = packagesToGenerate;

// TODO: Send the fully merged config to the instead of the temporary scoped config
const mergedConfig = merge(tutoneConfig, cfg)
const tutoneConfigYAML = yaml.stringify(cfg);

console.log('');
console.log('Tutone config:');
console.log('');
console.log(tutoneConfigYAML);
console.log('');

// Check to see which mutations the client is missing
const schemaMutations = schemaLatest.mutationType.fields.map(field => field.name);
const clientMutationsDiff = schemaMutations.filter(x => !clientMutations.includes(x));

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

const listOfPackagesToGenerate = packagesToGenerate.map(pkg => pkg.name);

console.log('List of packages to generate:', listOfPackagesToGenerate);

module.exports = {
  heroMention,
  schemaMutations,
  clientMutations,
  clientMutationsDiff,
  newApiMutationsMsg,
  clientMutationsDiffMsg,
  changedEndpoints,
  tutoneConfig: tutoneConfigYAML,
  packagesToGenerate: listOfPackagesToGenerate, // pop() is temporary for testing
};
