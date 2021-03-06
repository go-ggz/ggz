{
  test:: {
    kind: 'pipeline',
    name: 'testing',
    platform: {
      os: 'linux',
      arch: 'amd64',
    },
    steps: [
      {
        name: 'generate',
        image: 'golang:1.13',
        pull: 'always',
        commands: [
          'make generate',
        ],
        volumes: [
          {
            name: 'gopath',
            path: '/go',
          },
        ],
      },
      {
        name: 'vet',
        image: 'golang:1.13',
        pull: 'always',
        commands: [
          'make vet',
        ],
        volumes: [
          {
            name: 'gopath',
            path: '/go',
          },
        ],
      },
      {
        name: 'lint',
        image: 'golang:1.13',
        pull: 'always',
        commands: [
          'make lint',
        ],
        volumes: [
          {
            name: 'gopath',
            path: '/go',
          },
        ],
      },
      {
        name: 'misspell',
        image: 'golang:1.13',
        pull: 'always',
        commands: [
          'make misspell-check',
        ],
        volumes: [
          {
            name: 'gopath',
            path: '/go',
          },
        ],
      },
      {
        name: 'embedmd',
        image: 'golang:1.13',
        pull: 'always',
        commands: [
          'make embedmd',
        ],
        volumes: [
          {
            name: 'gopath',
            path: '/go',
          },
        ],
      },
      {
        name: 'test',
        image: 'golang:1.13',
        pull: 'always',
        commands: [
          'make test',
        ],
        volumes: [
          {
            name: 'gopath',
            path: '/go',
          },
        ],
      },
      {
        name: 'codecov',
        image: 'robertstettner/drone-codecov',
        pull: 'always',
        settings: {
          token: { 'from_secret': 'codecov_token' },
        },
      },
    ],
    volumes: [
      {
        name: 'gopath',
        temp: {},
      },
    ],
  },

  build(name, os='linux', arch='amd64', cgo=false)::
    local build_sqlite = if cgo then "-tags 'sqlite sqlite_unlock_notify'" else "";
    local build_static = if cgo then "-extldflags -static" else "";
    {
      kind: 'pipeline',
      name: name + '-' + os + '-' + arch,
      platform: {
        os: os,
        arch: arch,
      },
      steps: [
        {
          name: 'build-push',
          image: 'golang:1.13',
          pull: 'always',
          environment: {
            CGO_ENABLED: if cgo then "1" else "0",
          },
          commands: [
            'make generate',
            'go build -v '+ build_sqlite +' -ldflags "'+ build_static +' -X github.com/go-ggz/ggz/pkg/version.Version=${DRONE_COMMIT_SHA:0:8} -X github.com/go-ggz/ggz/pkg/version.BuildDate=`date -u +%Y-%m-%dT%H:%M:%SZ`" -a -o release/' + os + '/' + arch + '/' + name + ' ./cmd/' + name,
          ],
          when: {
            event: {
              exclude: [ 'tag' ],
            },
          },
        },
        {
          name: 'build-tag',
          image: 'golang:1.13',
          pull: 'always',
          environment: {
            CGO_ENABLED: if cgo then "1" else "0",
          },
          commands: [
            'make generate',
            'go build -v '+ build_sqlite +' -ldflags "'+ build_static +' -X github.com/go-ggz/ggz/pkg/version.Version=${DRONE_TAG##v} -X github.com/go-ggz/ggz/pkg/version.BuildDate=`date -u +%Y-%m-%dT%H:%M:%SZ`" -a -o release/' + os + '/' + arch + '/' + name + ' ./cmd/' + name,
          ],
          when: {
            event: [ 'tag' ],
          },
        },
        {
          name: 'executable',
          image: 'golang:1.13',
          pull: 'always',
          commands: [
            './release/' + os + '/' + arch + '/' + name + ' --help',
          ],
        },
        {
          name: 'dryrun',
          image: 'plugins/docker:' + os + '-' + arch,
          pull: 'always',
          settings: {
            daemon_off: false,
            dry_run: true,
            tags: os + '-' + arch,
            dockerfile: 'docker/' + name + '/Dockerfile.' + os + '.' + arch,
            repo: 'goggz/' + name,
            cache_from: 'goggz/' + name,
          },
          when: {
            event: [ 'pull_request' ],
          },
        },
        {
          name: 'publish',
          image: 'plugins/docker:' + os + '-' + arch,
          pull: 'always',
          settings: {
            daemon_off: 'false',
            auto_tag: true,
            auto_tag_suffix: os + '-' + arch,
            dockerfile: 'docker/' + name + '/Dockerfile.' + os + '.' + arch,
            repo: 'goggz/' + name,
            cache_from: 'goggz/' + name,
            username: { 'from_secret': 'docker_username' },
            password: { 'from_secret': 'docker_password' },
          },
          when: {
            event: {
              exclude: [ 'pull_request' ],
            },
          },
        },
      ],
      depends_on: [
        'testing',
      ],
      trigger: {
        ref: [
          'refs/heads/master',
          'refs/pull/**',
          'refs/tags/**',
        ],
      },
    },

  release:: {
    kind: 'pipeline',
    name: 'release-binary',
    platform: {
      os: 'linux',
      arch: 'amd64',
    },
    steps: [
      {
        name: 'generate',
        image: 'golang:1.13',
        pull: 'always',
        commands: [
          'make generate',
        ],
        volumes: [
          {
            name: 'gopath',
            path: '/go',
          },
        ],
      },
      {
        name: 'build-all-binary',
        image: 'golang:1.13',
        pull: 'always',
        commands: [
          'make release'
        ],
        when: {
          event: [ 'tag' ],
        },
        volumes: [
          {
            name: 'gopath',
            path: '/go',
          },
        ],
      },
      {
        name: 'deploy-all-binary',
        image: 'plugins/github-release',
        pull: 'always',
        settings: {
          files: [ 'dist/release/*' ],
          api_key: { 'from_secret': 'github_release_api_key' },
        },
        when: {
          event: [ 'tag' ],
        },
      },
    ],
    depends_on: [
      'testing',
    ],
    trigger: {
      ref: [
        'refs/tags/**',
      ],
    },
  },

  notifications(name, os='linux', arch='amd64', depends_on=[]):: {
    kind: 'pipeline',
    name: name + '-notifications',
    platform: {
      os: os,
      arch: arch,
    },
    steps: [
      {
        name: 'manifest',
        image: 'plugins/manifest',
        pull: 'always',
        settings: {
          username: { from_secret: 'docker_username' },
          password: { from_secret: 'docker_password' },
          spec: 'docker/' + name + '/manifest.tmpl',
          ignore_missing: true,
        },
      },
    ],
    depends_on: depends_on,
    trigger: {
      ref: [
        'refs/heads/master',
        'refs/tags/**',
      ],
    },
  },

  signature(key):: {
    kind: 'signature',
    hmac: key,
  }
}
