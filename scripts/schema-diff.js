module.exports = async ({
  core
}) => {
  const fs = require('fs');
  const yaml = require('yaml');

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

  // Check to see which mutations the client is missing
  const schemaMutations = schemaLatest.mutationType.fields.map(field => field.name);
  const clientMutationsDiff = schemaMutations.filter(x => !clientMutations.includes(x));

  console.log('Client Mutations:', clientMutations);
  console.log('Client is still missing the following mutations:\n', clientMutationsDiff);

  let newApiMutationsMsg = 'No new mutations since last check';
  if (endpointsDiff.length > 0) {
    heroMention = '@hero';
    newApiMutationsMsg = `'${endpointsDiff.join('\n')}'`;
  }

  let clientMutationsDiffMsg = ''
  if (clientMutationsDiff.length > 0) {
    clientMutationsDiffMsg = `'${clientMutationsDiff.join('\n')}'`;
  }


  core.setOutput('hero_mention', heroMention);
  core.setOutput('total_api_mutations_count', schemaMutations.length);
  core.setOutput('client_mutations_count', clientMutations.length);
  core.setOutput('client_mutations_missing_count', clientMutationsDiff.length);

  core.setOutput('new_api_mutations', newApiMutationsMsg);
  core.setOutput('client_mutations_missing', clientMutationsDiffMsg);

  await core.summary
  .addHeading('New Relic Client Go | NerdGraph API Report')
  .addRaw('Client mutations:')
  .addList(clientMutations)
  .addRaw('Client is missing the following mutations:')
  .addList(clientMutationsDiff)
  .write()
}
