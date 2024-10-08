# Changelog

> [!IMPORTANT]
> This project has not reached v1.0.0 and as such provides no backwards compatibility guarantees between versions.
> Pre v1.0.0 minor bumps will repesent breaking changes.

## v0.1.1 (2024-08-22)

* cmd/blackbox: fix active revision when listing release revisions ([e08dd7f](https://github.com/davidmdm/yoke/commit/e08dd7fb82cf2c611c20687770097c39aebda056))
* chore: update dependencies ([ec169ac](https://github.com/davidmdm/yoke/commit/ec169ac9ca9aab5f142ed3396ccc6b9e29168c85))
* yokecd: use yoke.EvalFlight instad of low-level wasi.Execute to be more compatible with pkg/Flight helpers ([87230e9](https://github.com/davidmdm/yoke/commit/87230e9e720c8e386c70ea1a86782408ec46f944))
* cmd/internal/changelog: add dates to tag ([6163eae](https://github.com/davidmdm/yoke/commit/6163eae045e3d5487d519414dc82b03337c5403a))
* cmd/internal/changelog: fix issue where multiple tags on same commit would only show one tag ([dae2c54](https://github.com/davidmdm/yoke/commit/dae2c543adc2bad74cb8ea62bfa9a539ce2791fc))
* cmd/internal/changelog: added internal command to generate changelog for project ([c98628b](https://github.com/davidmdm/yoke/commit/c98628b6443eed0029acf368e5ab12f57ad7c8ef))

## v0.1.0 (2024-06-22)

> [!CAUTION]
> This version contains breaking changes, and is not expected to be compatible with previous versions

* yoke: breaking change: represent revision history as multiple secrets ([cde0d83](https://github.com/davidmdm/yoke/commit/cde0d832f855f26ea51e6385677fdbd5f2d92e41))

## v0.0.11 (2024-06-17)

* yoke/takeoff: switch default to true for --create-crds flag ([4ffe721](https://github.com/davidmdm/yoke/commit/4ffe7218468e2b3a5897af5c2bfd42eca9439de9))
* cmd: added --poll flag to set poll interval for resource readiness during takeoff and descent ([63a6437](https://github.com/davidmdm/yoke/commit/63a64376c8ac32f564144eb9ece290fa9d992e6c))

## yokecd-installer/v0.0.6 (2024-06-16)

* yokecd-installer: bump argocd chart to version 7.1.3 ([ea662ae](https://github.com/davidmdm/yoke/commit/ea662ae7dbd55b3ac6605bbd325578151d265588))

## v0.0.10 (2024-06-15)

* deps: update project dependencies ([2785be6](https://github.com/davidmdm/yoke/commit/2785be63452ff98263ebca85dd74c1bc07bdecee))

## v0.0.9 (2024-06-15) - yokecd-installer/v0.0.5 (2024-06-15)

* pkg/helm: support subchart dependencies ([969e592](https://github.com/davidmdm/yoke/commit/969e592ef4b8555b30c84f380b0d4a362a05620c))
* cmd/takeoff: test --wait option ([14f3c67](https://github.com/davidmdm/yoke/commit/14f3c670f5508724f475e938d0db6f2d8e1fcd0d))
* pkg/yoke: add wait option to descent ([e7580be](https://github.com/davidmdm/yoke/commit/e7580be0ce06d9a536c8f4686b33f3728e8688a7))
* pkg/yoke: add wait option to takeoff ([2721de8](https://github.com/davidmdm/yoke/commit/2721de8060196723c78981e2896701e1671a7773))
* internal/k8s: support readiness checks for workloads resources like pods, deployments, statefulsets, replicasets, and so on ([21d2e7c](https://github.com/davidmdm/yoke/commit/21d2e7c623fe752ac2e2b23c03cdb6f442857afd))
* pkg/yoke: move wasm related code into same file ([2067753](https://github.com/davidmdm/yoke/commit/2067753d8e3ed607a03d54c0fc89c7ce3c1bf51e))
* yoke/debug: add debug timers to descent, mayday, and turbulence commands ([f377a27](https://github.com/davidmdm/yoke/commit/f377a27580d3ead62bb19504adaba348bc11c09c))
* yoke/takeoff: wait for namespace created by -namespace to be ready ([178bf8d](https://github.com/davidmdm/yoke/commit/178bf8d3e1c37d4d310c7baba0ab2a71890a8821))

## v0.0.8 (2024-06-01)

* pkg/yoke: set release in env of flight; update pkg/flight accordingly ([488985e](https://github.com/davidmdm/yoke/commit/488985e3aa36c7a579a5c220ddb30c17e754063d))

## v0.0.7 (2024-06-01)

* cmd/yoke: add create namespace and crd logic to takeoff (#20)

* cmd/yoke: add create namespace and crd logic to takeoff
* pkg/yoke: refactor move all takeoff command logic into commander.Takeoff ([5aebdcc](https://github.com/davidmdm/yoke/commit/5aebdccb99ccf63a595052b269598756c4d83faf))

## yokecd-installer/v0.0.4 (2024-05-29)

* pkg/helm: do not render test files ([df2329f](https://github.com/davidmdm/yoke/commit/df2329f7beb24366097c2a7547225304ebd766bf))
* yoke: use stdout to determine color defaults for takeoff and turbulence ([164c7b7](https://github.com/davidmdm/yoke/commit/164c7b79e06496092fb7b0d9114ef363910d3f38))
* yoke: concurrently apply resources during takeoff ([50cad15](https://github.com/davidmdm/yoke/commit/50cad159a12dee79d2461b10455bf3828151ffe4))
* yoke: rename global -verbose flag to -debug ([a9a803c](https://github.com/davidmdm/yoke/commit/a9a803c4d3dbb61250ea73b852dec6aeb6d6075a))

## v0.0.6 (2024-05-19)

* yoke: add takeoff diff-only tests ([824d4fb](https://github.com/davidmdm/yoke/commit/824d4fb75c4c6040695a1c4a5c414ead59ffb9f7))
* refactor: stdio, consolidate use of canonical object map ([f5e2dff](https://github.com/davidmdm/yoke/commit/f5e2dff4e09528d5c7f70f11b0d53ea72fecc950))
* formatting: fix import order ([124d8a6](https://github.com/davidmdm/yoke/commit/124d8a67adfe4f3fdb3d4e6f5367c719c014f0f3))
* refactor: add contextual stdio for better testing ([91c8391](https://github.com/davidmdm/yoke/commit/91c8391444f01216c7aadf45917d70b4148c8d70))
* yoke: update xcontext dependency ([9c6c178](https://github.com/davidmdm/yoke/commit/9c6c178243cfec6dd82d716b0f173c58f6967bf9))
* yoke: use canonical object map for takeoff diffs ([7a9f0ff](https://github.com/davidmdm/yoke/commit/7a9f0ffbc8d069be395cc1f88293e13796135d64))
* Added --diff-only to takeoff command (#17) ([e4c8a25](https://github.com/davidmdm/yoke/commit/e4c8a258e8d99cea04c9842fae6e51e30e042307))

## v0.0.5 (2024-05-18)

* yoke: drift detection ([3ab27a7](https://github.com/davidmdm/yoke/commit/3ab27a7610ab869830807bbe17cd51895e8f8a6b))
* yoke: add drift detection ([3e1e2a9](https://github.com/davidmdm/yoke/commit/3e1e2a98fbce95f5435e6e6f3fe1dbfd7bd87d22))
* readme: add link to official documentation ([bdf3565](https://github.com/davidmdm/yoke/commit/bdf3565f8e89abb6745fe5c2a1ffa6a2d14d1217))

## yokecd-installer/v0.0.3 (2024-05-04)

* yokecd-installer: make yokecd docker image version configurable ([821d6e3](https://github.com/davidmdm/yoke/commit/821d6e3ee992f7ff25b75ba3ea84d55a85bae5f5))

## v0.0.4 (2024-04-29)

* yoke: add namespace to debug timer ([4e8ab04](https://github.com/davidmdm/yoke/commit/4e8ab04a46382649e3b52930fb0f590b2fc3a5a2))
* refactor: fix import orderings ([6d3a09f](https://github.com/davidmdm/yoke/commit/6d3a09f3aed6e77fc42a9d8ff06289b91e51999a))
* yoke: ensure namespace exists before applying resources ([8cee965](https://github.com/davidmdm/yoke/commit/8cee96515043712dc2caca04aea7475fa78a2506))
* yoke: fix help text mistakes ([a6657ae](https://github.com/davidmdm/yoke/commit/a6657ae620d3d80c5f24c81c172e34d76b62c979))
* yokecd: remove wasm after use in build mode ([7dbd330](https://github.com/davidmdm/yoke/commit/7dbd330cea6c1d3fcf51390d3c4ab257968bb520))

## v0.0.3 (2024-04-25)

* yokecd: added config parsing tests ([7dc8200](https://github.com/davidmdm/yoke/commit/7dc8200c13f9750cc93f093586dcec227d883a25))
* yokecd: add build mode ([8760b9f](https://github.com/davidmdm/yoke/commit/8760b9f94723a161281192187bd09bc30ddfe499))

## yokecd-installer/v0.0.2 (2024-04-21)

* releaser: fix patch inference ([960853a](https://github.com/davidmdm/yoke/commit/960853a4ab113904077430db84081ac264685b4c))
* pkg: added flight package with convenience functions for flight execution contexts ([fd401ea](https://github.com/davidmdm/yoke/commit/fd401ea19d5bb304fb4b7b45245c55e9c689c615))
* yokecd: require wasm field at config load time ([660a913](https://github.com/davidmdm/yoke/commit/660a913f6ec1fae7fe216cb3a4b0c8dbb144d6a2))

## v0.0.2 (2024-04-20)

* yoke: added verbose mode with debug timings for functions ([2f87cef](https://github.com/davidmdm/yoke/commit/2f87cef5cf06e757f73dcc21643aa38117fe24c2))
* yoke: improve takeoff help text ([b74f17d](https://github.com/davidmdm/yoke/commit/b74f17d1599a6cb7ce49b202145fed3663a5dad7))
* yoke: add wazero to version output ([af90ae6](https://github.com/davidmdm/yoke/commit/af90ae6624b7687982f8f74f0ee890cc87e9ee41))

## v0.0.1 (2024-04-20) - yokecd-installer/v0.0.1 (2024-04-20)

* releaser: release patch versions from now on ([eac2db4](https://github.com/davidmdm/yoke/commit/eac2db4c409ab28039c4cddc86c1c4a96f380553))
* update dependencies ([44a6dd7](https://github.com/davidmdm/yoke/commit/44a6dd79af61344d345804feb894a843eedb6653))
* yoke: fix force conflicts flag not propagated ([a8a086c](https://github.com/davidmdm/yoke/commit/a8a086c210e04f2323b8f44289fdf138a9204186))
* yoke: interpret http path suffixes with .gz as gzipped content ([e68f8ba](https://github.com/davidmdm/yoke/commit/e68f8ba32e3b96bf7dc93a342799ece3a8f8623b))

## yokecd-installer/v0.0.0-20241704012137 (2024-04-16)

* yokecd: use argo-helm/argocd for installer ([a3fe4df](https://github.com/davidmdm/yoke/commit/a3fe4df441404ab2a9b1225175e8ed3c2fac603c))
* yoke: use secrets instead of configmaps for storing revision state ([5e39717](https://github.com/davidmdm/yoke/commit/5e397171e214463166f5facfa62050f0f60324fd))
* add tests to workflow ([01608f8](https://github.com/davidmdm/yoke/commit/01608f8a5df6186fbd5522f1440b4990133db177))
* revert wazero to v1.6.0 and use compiler ([c5d48bf](https://github.com/davidmdm/yoke/commit/c5d48bf0556b28df698c438747ed6d3d02a15e38))

## yokecd-installer/v0.0.0-20241004031222 (2024-04-09)

* fix compressor ([5dc59c0](https://github.com/davidmdm/yoke/commit/5dc59c0721aa56bf77ceb9995dd46b5e49688446))
* fix download err ([2552acb](https://github.com/davidmdm/yoke/commit/2552acbfc0a563475bc9013a8a54877535069840))
* release yokecd-installer ([ed4d68d](https://github.com/davidmdm/yoke/commit/ed4d68d16f5153d2b1e399006cfd4b8faaff581e))
* release yokecd as gz ([95e83db](https://github.com/davidmdm/yoke/commit/95e83db497814b20f47496394f017be5cf947ac8))
* add yokecd releaser workflow ([579ca0c](https://github.com/davidmdm/yoke/commit/579ca0c0452e9346f4dfd40899dd8ca4ed727916))
* yokecd: remove leading slashes for wasm parameter ([9a69c9e](https://github.com/davidmdm/yoke/commit/9a69c9e2a5a5ac99f609cc2ad0103e3ed6a51b6b))
* support one parameter wasm instead of wasmURL and wasmPath ([02245e2](https://github.com/davidmdm/yoke/commit/02245e28ee7257c4765e757c5ebd9a805b41e9a2))
* yoke: add resource ownership conflict test ([8809ad7](https://github.com/davidmdm/yoke/commit/8809ad730aceb2eaf8b7171bdf4f482199e85b11))
* yoke: support gzip wasms ([8d3dbb1](https://github.com/davidmdm/yoke/commit/8d3dbb144650d8a68910fded2033b21b5b868302))
* yokecd: add suport for wasmPath parameter ([cfc3952](https://github.com/davidmdm/yoke/commit/cfc39526323d6e295fc03ff98f299a81f90b2dba))
* test simplified yokecd application ([70c44d3](https://github.com/davidmdm/yoke/commit/70c44d3bf34b8aeb005eedcbdb564cace0279492))

## v0.0.0-beta3 (2024-03-24)

* hardcode yokecd as plugin name, and simplified plugin definition ([0de06a9](https://github.com/davidmdm/yoke/commit/0de06a959bde14d7432f737d60e9a77db88e79a1))
* removed yoke exec in favor of takeoff --test-run ([a88793b](https://github.com/davidmdm/yoke/commit/a88793b08995540b31bb3715b08e4cb084bacbf1))
* use wazero interpreter instead of compiler ([f87011a](https://github.com/davidmdm/yoke/commit/f87011a242d62d8a03658b6e55da1127cab7de70))
* fix http proto check for wasm loading ([b8cd522](https://github.com/davidmdm/yoke/commit/b8cd522c7f82bf69a2feb3838d47d86611d14358))
* added yoke exec for debugging wasm ([e77f607](https://github.com/davidmdm/yoke/commit/e77f6078b503753996b2ea6dc21cdad6ca210dd1))
* revert wazero to v1.6.0 ([5402e44](https://github.com/davidmdm/yoke/commit/5402e44fc4911acf191c25b885f6f90aae643ec9))
* updates deps ([e69bbfd](https://github.com/davidmdm/yoke/commit/e69bbfdea0c0fd4a6f779cfed4b0d035bf9d0295))
* make flight marshalling less verbose by omitting app source ([2eaa0db](https://github.com/davidmdm/yoke/commit/2eaa0db79666308c2a7b18487dda5dd25936c65d))
* remove website ([407617c](https://github.com/davidmdm/yoke/commit/407617c7a816bd6714cfa8ea469ddc22a3ff08d4))
* update debug logs ([da35b1e](https://github.com/davidmdm/yoke/commit/da35b1ea32eb2ccf7c97b742003b29d56f15a338))
* fix sync policy support ([0c99311](https://github.com/davidmdm/yoke/commit/0c993113f8563e0360079f1ded923fec05c5dca7))
* add basic syncPolicy support ([9361903](https://github.com/davidmdm/yoke/commit/9361903b8cdbff65d6e37e681ab6cde5d1f4a210))
* add plugin parameter ([602a6e7](https://github.com/davidmdm/yoke/commit/602a6e79eaf4e2abec50a21b4b442d812246ec78))
* fixed flight rendering logic in yokecd ([12351c1](https://github.com/davidmdm/yoke/commit/12351c190aae3030dc7f7b5bab7954be86b7e1a4))
* make flight spec embed application spec ([fd744d5](https://github.com/davidmdm/yoke/commit/fd744d52dd3bd7c57b46ffdb1511d49e3986030c))
* yokecd in progress ([541cdac](https://github.com/davidmdm/yoke/commit/541cdac548f54caf65faaa0ebcb15dfdf812bc51))
* add more debug info to yokecd ([34e269a](https://github.com/davidmdm/yoke/commit/34e269ae906afe24d96132f2524c801ff09c80f5))
* add yokecd example flight with patched argocd-repo-server ([d24024b](https://github.com/davidmdm/yoke/commit/d24024b123b2aa60cc21c0e2a32718c32572ae03))
* fix go.sum ([1327974](https://github.com/davidmdm/yoke/commit/13279741f03bceccc38e0481dbcede48bb497abd))
* basic code for yokecd ([396ccfc](https://github.com/davidmdm/yoke/commit/396ccfc4d2a4cb6002e8e273f583674af18f38f7))
* add version to helm2go ([08bfcac](https://github.com/davidmdm/yoke/commit/08bfcaca2c3a61f4733cc89aa7b303840b5970d8))
* update helm2go to default to map[string]any if cannot generate values type ([48b3a22](https://github.com/davidmdm/yoke/commit/48b3a22d16b9f4bbd788865b077fecc695488756))
* add force-conflicts flag for takeoff ([c058acb](https://github.com/davidmdm/yoke/commit/c058acb49961f7a67822a4e928cd069d237aa776))

## v0.0.0-beta2 (2024-03-15)

* use server-side apply patch ([68f1d97](https://github.com/davidmdm/yoke/commit/68f1d9716d0f29788aef1a831a9d958a94bcc98d))
* use official argocd install yaml for argo example ([bb46eb3](https://github.com/davidmdm/yoke/commit/bb46eb3d25ef9fe2492884a311ce83fbb595c35c))
* try create before apply ([172fc7f](https://github.com/davidmdm/yoke/commit/172fc7f758e58afeacf4138d9fba247551d85149))
* support non-namespaced resources ([abcd57a](https://github.com/davidmdm/yoke/commit/abcd57a4feab395111d54e11ced9d9d36acc3dd7))
* add skipDryRun and namespace flags to takeoff ([86d2081](https://github.com/davidmdm/yoke/commit/86d2081017d3645a73afa93885458dddd24e5a74))
* minor refactor of k8s client ([1c5f65a](https://github.com/davidmdm/yoke/commit/1c5f65a2c9a126209cfc9a9a96644748e5f2477e))
* udpated helm2go output and added mongodb example ([7227a24](https://github.com/davidmdm/yoke/commit/7227a24cae80035f9928004f84f68c4dc41f0771))
* add schema bool flag to helm2go ([cd9b074](https://github.com/davidmdm/yoke/commit/cd9b074220060e1819432a9d29bdccaba0f6a927))
* helm2go: infer from values ([a00cb9c](https://github.com/davidmdm/yoke/commit/a00cb9cc68f9582e61b67e7686e5e77556667e65))
* redis example uses generated flight from helm2go ([ce6a82e](https://github.com/davidmdm/yoke/commit/ce6a82e6e6eb9dd9fbf9088a24dffc7263976552))
* helm2go generates flight package instead of type file ([b123bb1](https://github.com/davidmdm/yoke/commit/b123bb1331905f61a554ac622629b84664784265))
* refactored helm2go ([5326c47](https://github.com/davidmdm/yoke/commit/5326c47323a8057d0a3aaf5b623cb986b0ea95b7))
* generated pg values.go using helm2go ([0c360bb](https://github.com/davidmdm/yoke/commit/0c360bb9cd7c1665b5150ecc7a4746f6029b9544))
* renamed platters to flights and added helm2go script ([9da0265](https://github.com/davidmdm/yoke/commit/9da02655f5e3985e18c3343ff64e7b589bd83735))

## v0.0.0-beta1 (2024-02-29)

* starting website ([dd8c995](https://github.com/davidmdm/yoke/commit/dd8c99584130ab84915a6bf5cc7e5c36af8de2a1))
* added apply dry run failure test ([10c65b7](https://github.com/davidmdm/yoke/commit/10c65b7b86d5f7976e4a6ff3e47ec262c8d50748))
* remove .DS_Store ([f52068a](https://github.com/davidmdm/yoke/commit/f52068a3ca6bf6fcfefb8631284f86718fb994c3))
* fix go.sum ([ae126b4](https://github.com/davidmdm/yoke/commit/ae126b4f6d917ecc8d962966831d98a859230c76))
* refactored tests to not use TestMain ([7a5213f](https://github.com/davidmdm/yoke/commit/7a5213f6be100fc50c65e5efd9f6a2f658c62a39))
* formatting ([607b346](https://github.com/davidmdm/yoke/commit/607b3462b5523fc18205cabadf1c0de40c043229))
* remove .DS_Store ([1400e5b](https://github.com/davidmdm/yoke/commit/1400e5b0f419cf6e1a670d6b5a0362b884261ada))
* the great renaming to yoke ([578ac2c](https://github.com/davidmdm/yoke/commit/578ac2cef7070fc234ab83058b23eab72248ef5a))
* ported descent and mayday to sdk ([a44cf26](https://github.com/davidmdm/yoke/commit/a44cf26e154ab52b85ad19cba6e057fa3547859c))
* started sdk restructuring ([6a58b9b](https://github.com/davidmdm/yoke/commit/6a58b9bc83baa99734915b30cc04e2f932f2566c))
* add export to stdout ([44b071f](https://github.com/davidmdm/yoke/commit/44b071f4bd25c489fa900f8bbda166039ea0ae2f))
* rename k8 to k8s ([3223e8b](https://github.com/davidmdm/yoke/commit/3223e8b3f81f43e22a161efa978c754fc9c04ed4))
* refactor kube example ([a5f85c4](https://github.com/davidmdm/yoke/commit/a5f85c4dbd6aa9d8d4d7ec00de759e7ffb474a4e))
* refactored example platters around ([30067fc](https://github.com/davidmdm/yoke/commit/30067fcf220410e5e1ea808908d4f143dc32b93c))
* wrote first test ([7f4e9a9](https://github.com/davidmdm/yoke/commit/7f4e9a9150f8173a34a3b498475155f3a389addf))
* add blackbox --mapping flag ([8506be6](https://github.com/davidmdm/yoke/commit/8506be61b987dd101265647ebbae98680ace479f))
* use all prefix for embedding private templates in helm expanded example ([f1850bf](https://github.com/davidmdm/yoke/commit/f1850bfd99bdde4a74cb788c913290586396f4ad))
* load helm chart from embed.FS work in progress ([0ce494e](https://github.com/davidmdm/yoke/commit/0ce494e2e736958823ddca201263c968ede65b58))
* added load chart from FS: wip ([70b3cee](https://github.com/davidmdm/yoke/commit/70b3cee1b6576df1ee5d14160b8e7046f6991621))
* update helmchart example to make it configurable ([c5133ef](https://github.com/davidmdm/yoke/commit/c5133ef509b318a3fd47eaafc8426b0e7ce0d844))
* update halloumi metadata ([11c9b2e](https://github.com/davidmdm/yoke/commit/11c9b2e654f2eb486af4e5fcf61d191ad5937771))

## v0.0.0-alpha17 (2024-02-25)

* updated helm api ([964d147](https://github.com/davidmdm/yoke/commit/964d147b1171920142533a87ce3868a23e2dccd1))
* initial support for helm chart compatibility ([d3c926e](https://github.com/davidmdm/yoke/commit/d3c926e94022635bab35f32d91108df675e1d7e5))

## v0.0.0-alpha16 (2024-02-24)

* update verison command to show k8 client-go version as well ([831fdd7](https://github.com/davidmdm/yoke/commit/831fdd7d5573cb1bda1b7c4f28d500d1403bec79))
* change diff indentation ([f3173be](https://github.com/davidmdm/yoke/commit/f3173be28d44754138e929b83245d0b103538970))

## v0.0.0-alpha15 (2024-02-24)

* print diff between revisions ([706e050](https://github.com/davidmdm/yoke/commit/706e0501a4b9789fae2801249ef7da9fe0cb3187))
* refactored revision source ([6aa96a5](https://github.com/davidmdm/yoke/commit/6aa96a5c38a54559ee93a7990f53b68bf0a0ccfa))
* added shas to revisions ([ce5a7da](https://github.com/davidmdm/yoke/commit/ce5a7da3531078c527eedcd5e09522f5a800e1b3))
* refactor blackbox ([fc4ad5a](https://github.com/davidmdm/yoke/commit/fc4ad5a722f885309e9054b37a3ecaf5c1d66cbf))
* update blackbox output ([505e281](https://github.com/davidmdm/yoke/commit/505e281d35810e3b80966d06c948dd3e210626bf))

## v0.0.0-alpha14 (2024-02-24)

* added mayday command ([d982624](https://github.com/davidmdm/yoke/commit/d982624ea20bb8dfd6b7702a13e96717797e507e))
* remove unnecessary newline from error ([2702985](https://github.com/davidmdm/yoke/commit/27029855fe5e16f1067c89c29fff4727881953d3))

## v0.0.0-alpha13 (2024-02-23)

* finish first pass at blackbox command ([558273b](https://github.com/davidmdm/yoke/commit/558273b069e75547c672ba19e881cd25b7b16c6d))
* update deps and formatting ([05b5096](https://github.com/davidmdm/yoke/commit/05b5096d29158e5ba26e701da4222360e978ceec))
* blackbox under construction ([91d3fa7](https://github.com/davidmdm/yoke/commit/91d3fa78c37795b21e3b0b626a1ea0c5393ea647))
* removed resource utility package in favor of applyconfigurations ([2de98b0](https://github.com/davidmdm/yoke/commit/2de98b0ed2679917f0f7cec389df3596538caaff))

## v0.0.0-alpha12 (2024-02-23)

* create an ownership check ([8c2d7f9](https://github.com/davidmdm/yoke/commit/8c2d7f9f4993e0e1c3140a991f8206ce9adca570))
* added blackbox shell ([3c34f8c](https://github.com/davidmdm/yoke/commit/3c34f8c6c44a8ba6acd480fbfe39640ed46fec45))

## v0.0.0-alpha11 (2024-02-21)

* first working pass of descent command ([91cc860](https://github.com/davidmdm/yoke/commit/91cc86088b8f82c3939089398731fb4480284581))
* first pass at descent command ([c71368b](https://github.com/davidmdm/yoke/commit/c71368bbdc8e953674e7f5f8533ea639754e8424))
* modified configmap structure ([cb0691f](https://github.com/davidmdm/yoke/commit/cb0691fc4ae2363f0dea13c0595f806c5f38286e))
* dynamic platter example ([043358a](https://github.com/davidmdm/yoke/commit/043358a269737ad7dfd53e302b7c2c4dd92d705f))

## v0.0.0-alpha10 (2024-02-19)

* updated canonical name to include api version and changed deploy to apply ([78980fb](https://github.com/davidmdm/yoke/commit/78980fb1b7f434d32417bf0ac33c9316faf4dcc4))
* adding to resource utility package ([4bf13be](https://github.com/davidmdm/yoke/commit/4bf13be489cdef9861ac2205002084e8aeeb1d55))

## v0.0.0-alpha9 (2024-02-18)

* do not apply identical revisisions but do a noop ([750de31](https://github.com/davidmdm/yoke/commit/750de31c9907f162a574bb5bf08b803a4da2e3a6))
* added beginning of a basic utility package for resource definitions ([cd43c11](https://github.com/davidmdm/yoke/commit/cd43c1173f568dbabfe10015c79cf70f73e5dc82))

## v0.0.0-alpha8 (2024-02-18)

* allow wasm executable to receive stdin as input ([d7d9922](https://github.com/davidmdm/yoke/commit/d7d992296e40451d9ea596976cbd27be007301dc))

## v0.0.0-alpha7 (2024-02-18)

* add outdir option to takeoff instead of render or runway command ([c860e31](https://github.com/davidmdm/yoke/commit/c860e315edbe951695d05ac238a1f6ffa5f860f8))

## v0.0.0-alpha6 (2024-02-18)

* support yaml encodings of platters ([8f4cde7](https://github.com/davidmdm/yoke/commit/8f4cde72f2cd506ebdf9e18ba09e8c49b041c86d))

## v0.0.0-alpha5 (2024-02-17)

* add single or multi resource platter support and stdin source support ([20fc25f](https://github.com/davidmdm/yoke/commit/20fc25fcb849627707e63edf0c6b8fc0213e75bb))
* fix newline after root help text if no command provided ([2d4e04f](https://github.com/davidmdm/yoke/commit/2d4e04fdd0dd23998896de8899cd3f43d4c16654))

## v0.0.0-alpha4 (2024-02-17)

* small refactoring ([acc3351](https://github.com/davidmdm/yoke/commit/acc335177e86337209fc2a2746398df8d03871be))
* add dry run before applying resources ([3cb6d54](https://github.com/davidmdm/yoke/commit/3cb6d54736c4d80c3331913bc9c95b50e2dea8aa))
* add halloumi logo to readme ([f98fb12](https://github.com/davidmdm/yoke/commit/f98fb129e683b8b1fbf7ae72ce4c00a65ee69b5b))
* update readme ([087d62c](https://github.com/davidmdm/yoke/commit/087d62c6e6adb6ecae29d3c26f24bebdb3079332))

## v0.0.0-alpha3 (2024-02-17)

* added readme, license, and more aliases ([a8e6152](https://github.com/davidmdm/yoke/commit/a8e615276d55e64debe3a73048c8ad12974f37d6))

## v0.0.0-alpha2 (2024-02-17)

* go directive 1.22 ([b965e58](https://github.com/davidmdm/yoke/commit/b965e584bae7252d4c87f14a5ad87a0df7642c27))

## v0.0.0-alpha1 (2024-02-17)

* added root command help text ([26c24a3](https://github.com/davidmdm/yoke/commit/26c24a36f09a724187c6a2fcfef913f67aefaee4))
* takeoff help text ([11a8602](https://github.com/davidmdm/yoke/commit/11a86027ec3cee286869a67e9ddd870c8caac352))
* formmatting ([94f9345](https://github.com/davidmdm/yoke/commit/94f9345f4a30ff9808239dcc3dd39cc359046df2))
* remove wasibuild go utility and replace with task file ([61ca598](https://github.com/davidmdm/yoke/commit/61ca59846263e3aa670d6d7f9ce569b586ea4593))
* refactored into subcommands ([2db7d5d](https://github.com/davidmdm/yoke/commit/2db7d5dc2a234b94812f67508d3c10a473f5ea7b))
* remove orphans ([7efdb5d](https://github.com/davidmdm/yoke/commit/7efdb5dfb05c4d2bad5e87cbd0f2c952af7e8f33))
* save revisions as unstructured resources ([51bab9b](https://github.com/davidmdm/yoke/commit/51bab9b8b1d4d75d4f1ba755f45ccb4cc419e8c5))
* add halloumi-release-label ([51d98f7](https://github.com/davidmdm/yoke/commit/51d98f78744c659862a920e068e08e3b4f2a7c80))
* add revisions ([55bf01e](https://github.com/davidmdm/yoke/commit/55bf01ed5db36c2ed0020790d6c0be6988d1ec54))
* update deps ([2a6072b](https://github.com/davidmdm/yoke/commit/2a6072b420cc4de1a5e676a9337feffb49ac405a))
* namespace support ([d4aee29](https://github.com/davidmdm/yoke/commit/d4aee296cd6a9016596effd6bd0d9403bb7152ad))
* basic annotations ([405dd75](https://github.com/davidmdm/yoke/commit/405dd75d224b4cd0d96802020d5274b0820e60bb))
* k8 successful basic apply ([33515ae](https://github.com/davidmdm/yoke/commit/33515aeafa46ff2c5bfb272456203123cb607d36))
* k8: wip ([e846862](https://github.com/davidmdm/yoke/commit/e8468620e6da606f347b8f4814bb2b6cc7ab1190))
* refactor ([5a59247](https://github.com/davidmdm/yoke/commit/5a59247f3787bb7fb7e55d9a1c110e9735482c14))
* allow haloumi packages to be invoked with flags ([a368fce](https://github.com/davidmdm/yoke/commit/a368fcecbe6c4898d2c403c08de347974ba3f006))
* add .gitignore ([bee5929](https://github.com/davidmdm/yoke/commit/bee5929fd62250548b231fe805180152fe0a8368))
* refactor ([379d1fb](https://github.com/davidmdm/yoke/commit/379d1fbbf2e4f365ba667820c649526e16ced209))
* small utility for building wasi ([8bc22d5](https://github.com/davidmdm/yoke/commit/8bc22d586bfd14229f68fb334401b2108ef8884a))
* first haloumi binary working ([08bfa45](https://github.com/davidmdm/yoke/commit/08bfa45dca5f21d1d8875962f1531838337e96da))
* starting haloumi ([47b28fc](https://github.com/davidmdm/yoke/commit/47b28fcfc3766575eab80a1c9e640bc33d5ffa28))
* initial wazero runtime ([ac081a8](https://github.com/davidmdm/yoke/commit/ac081a89136e9e57abb27ac3797fc72095db6af9))

