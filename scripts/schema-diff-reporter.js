module.exports = async ({
  core
}) => {
  const diff = require('./schema-differ');

  core.setOutput('hero_mention', diff.heroMention);
  core.setOutput('total_api_mutations_count', diff.schemaMutations.length);
  core.setOutput('client_mutations_count', diff.clientMutations.length);
  core.setOutput('client_mutations_missing_count', diff.clientMutationsDiff.length);

  core.setOutput('new_api_mutations', diff.newApiMutationsMsg);
  core.setOutput('client_mutations_missing', diff.clientMutationsDiffMsg);
  core.setOutput('tutone_config', diff.tutoneConfig);
  core.setOutput('packages', diff.packagesToGenerate.toString());

  const gitHubActionSummary = core.summary

  gitHubActionSummary.addHeading('New Relic Client Go | NerdGraph API Report');


  if (diff.packagesToGenerate.length) {
    gitHubActionSummary
      .addRaw('Tutone will generate code for the following packages:')
      .addList(diff.packagesToGenerate);
  }

  await gitHubActionSummary
    .addRaw('Client mutations:')
    .addList(diff.clientMutations)
    .addRaw('Client is missing the following mutations:')
    .addList(diff.clientMutationsDiff)
    .write()
}

