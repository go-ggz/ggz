local pipeline = import 'pipeline.libsonnet';
local ggzServer = 'ggz-server';
local ggzRedirect = 'ggz-redirect';
[
  pipeline.test,
  pipeline.build(ggzServer, 'linux', 'amd64', true),
  pipeline.build(ggzServer, 'linux', 'arm64', true),
  pipeline.build(ggzServer, 'linux', 'arm', true),
  pipeline.build(ggzRedirect, 'linux', 'amd64'),
  pipeline.build(ggzRedirect, 'linux', 'arm64'),
  pipeline.build(ggzRedirect, 'linux', 'arm'),
  pipeline.release,
  pipeline.notifications(ggzServer, depends_on=[
    ggzServer + '-linux-amd64',
    ggzServer + '-linux-arm64',
    ggzServer + '-linux-arm',
    'release-binary',
  ]),
  pipeline.notifications(ggzRedirect, depends_on=[
    ggzRedirect + '-linux-amd64',
    ggzRedirect + '-linux-arm64',
    ggzRedirect + '-linux-arm',
    'release-binary',
  ]),
]
