// module.exports = async ({
//   core
// }) => {

const fs = require('fs');
const yaml = require('yaml');

let tutoneConfig = null
let schema = null

try {
  const tutoneConfigFile = fs.readFileSync('.tutone.yml', 'utf8')
  tutoneConfig = yaml.parse(tutoneConfigFile)

  const schemaFileLatest = fs.readFileSync('schema.json', 'utf8');
  schema = JSON.parse(schemaFileLatest);
} catch(err) {
  console.error(err);
}

const package = tutoneConfig.packages.find(x => x.name === 'entities');
// const packageTypesByName = package.types.filter(t => !t.skip_type_create && t.name.includes("Entity")).map(t => t.name);
const schemaTypes = schema.types.filter(t => t.name.endsWith("Entity") || t.name.endsWith("EntityOutline")).map(t => t.name)

// const typesDiff = schemaTypes.filter(x => {
//   console.log("schema type:", x);
//   console.log("pkg type:   ", packageTypesByName.includes(x));

//   return !packageTypesByName.includes(x)
// });

// console.log('packageTypes:', JSON.stringify(packageTypesByName, null, 2));
// console.log('schemaTypes: ', JSON.stringify(schemaTypes, null, 2));
// console.log('typesDiff:   ', typesDiff);
// console.log('typesDiff:   ', typesDiff.length);


const packageTypes = getPackageTypesFromTypesFile('/Users/sblue/dev/newrelic-client-go/pkg/entities/types.go');
const packageTypesLatest = getPackageTypesFromTypesFile('/Users/sblue/Desktop/2024.06.04-entities-package-types-latest.json');

function getPackageTypesFromTypesFile(path) {
  let matches = [];

  // Get get text between `type` and `struct` from the types.go file.
  const regex = /type\s+(\S+)\s+struct\s*{/g;

  try {
    // const text = fs.readFileSync(path, 'utf8');
    // return JSON.parse(text);

    const text = fs.readFileSync(path, 'utf8');

    let match;
    while ((match = regex.exec(text)) !== null) {
      matches.push(match[1]);
    }

    return matches;
  } catch (error) {
    console.error(error);
  }
}

const typesDiff = packageTypes.filter(x => {
  // console.log("pkg type:      ", x);
  // console.log("is latest type:", packageTypesLatest.includes(x));

  return !packageTypesLatest.includes(x)
});

console.log('typesDiff: ', JSON.stringify(typesDiff, null, 2));
console.log('typesDiff: ', typesDiff.length);
