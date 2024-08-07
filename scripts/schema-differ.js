const fs = require('fs');
const yaml = require('yaml');

// This name must match the alias we created in the #oac-automation-reports Slack channel.
const heroAliasName = '@oac-automation-watchers';

let tutoneConfig = null;
let schemaOld = null;
let schemaLatest = null;
let heroMention = "";

try {
  const tutoneConfigFile = fs.readFileSync('.tutone.yml', 'utf8')
  tutoneConfig = yaml.parse(tutoneConfigFile)

  const schemaFileOld = fs.readFileSync('schema-old.json', 'utf8');
  schemaOld = JSON.parse(schemaFileOld);

  const schemaFileLatest = fs.readFileSync('schema.json', 'utf8');
  schemaLatest = JSON.parse(schemaFileLatest);
} catch (err) {
  console.error(err);
}

// Check for any newly added mutations
const endpointsOld = schemaOld.mutationType.fields.map(field => field.name);
const endpointsLatest = schemaLatest.mutationType.fields.map(field => field.name);
const endpointsDiff = endpointsLatest.filter(x => !endpointsOld.includes(x));

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

console.log('');
console.log('New changedEndpoints', JSON.stringify(changedEndpoints, null, 2));
console.log('');

// Check to see which mutations the client is missing
const schemaMutations = schemaLatest.mutationType.fields.map(field => field.name);
const clientMutationsDiff = schemaMutations.filter(x => !clientMutations.includes(x));

console.log('');
console.log('Client is still missing the following API mutations:\n', clientMutationsDiff);
console.log('');

let newApiMutationsMsg = 'No new mutations since last check';
if (endpointsDiff.length > 0) {
  heroMention = heroAliasName;
  newApiMutationsMsg = `'${endpointsDiff.join('\n')}'`;
}

let clientMutationsDiffMsg = ''
if (clientMutationsDiff.length > 0) {
  clientMutationsDiffMsg = `'${clientMutationsDiff.join('\n')}'`;
}

function compareObjects(obj1, obj2, path = '') {
  let differences = [];

  for (let key in obj1) {
    if (obj1.hasOwnProperty(key) && obj2.hasOwnProperty(key)) {
      if (typeof obj1[key] === 'object' && typeof obj2[key] === 'object') {
        differences = differences.concat(compareObjects(obj1[key], obj2[key], `${path}.${key}`));
      } else if (obj1[key] !== obj2[key]) {
        differences.push({
          old: obj1[key],
          new: obj2[key],
          property: `${path}.${key}`,
        });
      }
    } else {
      differences.push({
        old: obj1[key],
        new: obj2[key] ? obj2[key] : undefined,
        property: `${path}.${key}`,
      });
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

module.exports = {
  heroMention,
  schemaMutations,
  clientMutations,
  clientMutationsDiff,
  newApiMutationsMsg,
  clientMutationsDiffMsg,
  changedEndpoints,
};
