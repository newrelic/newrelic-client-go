module.exports = async ({
  github,
  context,
  core
}) => {
  const owner = context.repo.owner;
  const repo = context.repo.repo;

  const workflows = await github.rest.actions.listRepoWorkflows({
    owner,
    repo
  });

  const workflow = workflows.data.workflows.find(w => w.path.includes(process.env.WORKFLOW_FILENAME));

  if (!workflow) {
    core.setFailed("No workflow found");
    return;
  }

  const runs = await github.rest.actions.listWorkflowRuns({
    owner,
    repo,
    workflow_id: workflow.id,
    status: "success",
    per_page: 10
  });

  let artifacts
  for (let i = 0; i < runs.data.workflow_runs.length; i++) {
    const result = await github.rest.actions.listWorkflowRunArtifacts({
      owner,
      repo,
      run_id: runs.data.workflow_runs[i].id
    });

    if (result.data.artifacts.length) {
      artifacts = result.data.artifacts
      break
    }
  }

  console.log("Artifacts:", JSON.stringify(artifacts, null, 2));

  if (runs.data.total_count === 0) {
    core.setFailed("No runs found");
    return;
  }

  const artifact = artifacts.find(artifact => artifact.name === process.env.ARTIFACT_NAME);
  if (artifact) {
    const response = await github.rest.actions.downloadArtifact({
      owner,
      repo,
      artifact_id: artifact.id,
      archive_format: 'zip'
    });
    require('fs').writeFileSync(process.env.ARTIFACT_FILENAME, Buffer.from(response.data));
    require('child_process').execSync(`unzip -o ${process.env.ARTIFACT_FILENAME} -d ${process.env.UNZIP_DIR}`);

    console.log("Artifact downloaded successfully");
  } else {
    core.setFailed("No artifact found");
  }
}
