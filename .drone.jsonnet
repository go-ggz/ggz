local pipeline = import 'pipeline.libsonnet';
local ggzServer = 'ggz-server';
local ggzRedirect = 'ggz-redirect';
[
  pipeline.test,
  pipeline.build(ggzServer, 'linux', 'amd64'),
  pipeline.build(ggzServer, 'linux', 'arm64'),
  pipeline.build(ggzServer, 'linux', 'arm'),
  pipeline.build(ggzRedirect, 'linux', 'amd64'),
  pipeline.build(ggzRedirect, 'linux', 'arm64'),
  pipeline.build(ggzRedirect, 'linux', 'arm'),
  pipeline.release,
  pipeline.notifications(depends_on=[
    ggzServer + '-linux-amd64',
    ggzServer + '-linux-arm64',
    ggzServer + '-linux-arm',
    ggzRedirect + '-linux-amd64',
    ggzRedirect + '-linux-arm64',
    ggzRedirect + '-linux-arm',
    'release-binary',
  ]),
]
