# Changelog

## [5.0.0](https://github.com/alvarolucio2007/Scholarly/compare/v4.2.0...v5.0.0) (2026-10-27)


### ⚠ BREAKING CHANGES

* Adding ErrTeacherNotFound if FK violation in UpdateCourse
* adding minimum password length to be 8
* CHanging writeJSONData to writeJSON in ListStudent (will get rid

### Features

* Adding (incomplete and buggy) user edit page ([1d72592](https://github.com/alvarolucio2007/Scholarly/commit/1d725927ff2b4807082a59b00be44123fe75b781))
* Adding basic dashboard page ([13db66d](https://github.com/alvarolucio2007/Scholarly/commit/13db66da588863d34962156438abb4660d6e5a46))
* Adding course dashboard, creation,edit and enrollment pages ([28925a8](https://github.com/alvarolucio2007/Scholarly/commit/28925a8a1480fc27368c49a169b2ee06d406e31b))
* Adding create student page ([02aa9b0](https://github.com/alvarolucio2007/Scholarly/commit/02aa9b06e2fa8f03ba38506743802800a3d3ea06))
* Adding enrollment dashboard and edit page ([96cdf60](https://github.com/alvarolucio2007/Scholarly/commit/96cdf6009ee7fc58a3fa73a45094f0d0996529d4))
* Adding ErrTeacherNotFound if FK violation in UpdateCourse ([068ec6c](https://github.com/alvarolucio2007/Scholarly/commit/068ec6c77b4a2e35e0c99a264d1eb75dd9f1e08d))
* Adding grades dashboard, editing, averages dashboard and average ([4ebe303](https://github.com/alvarolucio2007/Scholarly/commit/4ebe303ad1aecd4c94e2275d3ba451ac06a55e34))
* Adding student edit page ([fb1e0af](https://github.com/alvarolucio2007/Scholarly/commit/fb1e0af8cf2de4e29587766014d10cfe7397ae5b))
* Adding student page ([b5e2f1e](https://github.com/alvarolucio2007/Scholarly/commit/b5e2f1e54ada4403c99212bb9eb92e9798943910))
* Adding teacher creation page ([4c6d3ff](https://github.com/alvarolucio2007/Scholarly/commit/4c6d3ff3e2640ce66b44163a11bc683f012d0369))
* Adding teacher dashboard ([ec398c0](https://github.com/alvarolucio2007/Scholarly/commit/ec398c0eae44548ab92c866303220371619baa70))
* Adding teacher editing page ([129eac3](https://github.com/alvarolucio2007/Scholarly/commit/129eac35a6fc6139ca0b6a0ed8914326fe5d83fd))
* Adding test creation,editing and dashboard pages ([f5f2ac6](https://github.com/alvarolucio2007/Scholarly/commit/f5f2ac62917401116eadc3aca5079c1a753945a6))
* Adding user dashboard page ([38e32aa](https://github.com/alvarolucio2007/Scholarly/commit/38e32aacde36ab39a7812fcd54e84f25d73d6e26))
* Making FK violation in Course creation return Teacher Not Found ([b1bd273](https://github.com/alvarolucio2007/Scholarly/commit/b1bd273d2158bde871dbf7cdc75b0adf4f9b9acb))


### Bug Fixes

* adding logger and replcaing teachers with students in List ([2c04ed6](https://github.com/alvarolucio2007/Scholarly/commit/2c04ed6ff41d483bbd1a8be1c246d1829177581d))
* adding minimum password length to be 8 ([beceb00](https://github.com/alvarolucio2007/Scholarly/commit/beceb00bef0a75b57f6cc4f81bc040612feea90f))
* Adding proper FK violation handling in Create ([3bfd653](https://github.com/alvarolucio2007/Scholarly/commit/3bfd6531c9f82f46c0bbad33ad4830ac696e9c39))
* CHanging writeJSONData to writeJSON in ListStudent (will get rid ([b49a593](https://github.com/alvarolucio2007/Scholarly/commit/b49a593f6c43d33a404278e28366cadae07d9fa2))
* Fixing user.email not updating proper emails ([e4d051c](https://github.com/alvarolucio2007/Scholarly/commit/e4d051c0647611ce5ef9fb2cb6677a1d1860003a))
* Fixing writeJSON to writeJSONData in List ([8f1f563](https://github.com/alvarolucio2007/Scholarly/commit/8f1f5634d9b1e265d6f96ade0826129116efcf10))
* removing trailing space behind min in Password ([133a4b5](https://github.com/alvarolucio2007/Scholarly/commit/133a4b5ebbc83ac67d0b6e7862441be4c21b441c))

## [4.2.0](https://github.com/alvarolucio2007/Scholarly/compare/v4.1.0...v4.2.0) (2026-10-26)


### Features

* Adding basic login page. ([10dcfc3](https://github.com/alvarolucio2007/Scholarly/commit/10dcfc3ee7d3782f8d85c88df6bb8933598b0400))
* Adding user registration page ([f1dcf6d](https://github.com/alvarolucio2007/Scholarly/commit/f1dcf6dd18996ccd79ba8a1a3d847773fce88da5))


### Bug Fixes

* Adding for tags in label ([f4ee661](https://github.com/alvarolucio2007/Scholarly/commit/f4ee6619c09870bd68da4b464c3ba8f780258a68))
* Adding success return message and proper label to input. ([3f1d1a2](https://github.com/alvarolucio2007/Scholarly/commit/3f1d1a2dc4ba2869d35d504af92a2bce27cd3b2c))
* Adding successMessage and errorMessage reset ([10a2242](https://github.com/alvarolucio2007/Scholarly/commit/10a2242e78ba63530c13918855dc2effefe337d1))
* changing resolve from /register to /login ([88d84d1](https://github.com/alvarolucio2007/Scholarly/commit/88d84d1d476cd39c2902c29e3fbff149df0a508e))
* fixing error return ([9c78dfe](https://github.com/alvarolucio2007/Scholarly/commit/9c78dfe6f45d86c33ae1de7327a4cf6cb54b483d))
* Refactoring error names so they're both passing and more clear. ([57a4838](https://github.com/alvarolucio2007/Scholarly/commit/57a483802eb750b02f7874b78ba3ae16a0ca78a3))

## [4.1.0](https://github.com/alvarolucio2007/Scholarly/compare/v4.0.0...v4.1.0) (2026-10-24)


### Features

* Implementing rate limiter (100rpm limit) ([59b075e](https://github.com/alvarolucio2007/Scholarly/commit/59b075e970edde5a16f875e0a13f940aea404826))

## [4.0.0](https://github.com/alvarolucio2007/Scholarly/compare/v3.0.1...v4.0.0) (2026-10-24)


### ⚠ BREAKING CHANGES

* adding user registration function (just CreateUser in a different
* replacing Route with Group, and adding Authentication route
* adding basic auth middleware to debug,metrics and swagger routes,

### Features

* adding basic auth middleware to debug,metrics and swagger routes, ([176f071](https://github.com/alvarolucio2007/Scholarly/commit/176f0711d53ffbd6cadb8fa0de50637345312a14))
* Adding PASETO token creation handler ([5f6c105](https://github.com/alvarolucio2007/Scholarly/commit/5f6c105570d374eaa90b9005a2266bdce765d33c))
* adding user registration function (just CreateUser in a different ([502404f](https://github.com/alvarolucio2007/Scholarly/commit/502404f29a2e50d2b5ea3b637188dc5a190f06e9))


### Bug Fixes

* Adding creds length verification ([44d730c](https://github.com/alvarolucio2007/Scholarly/commit/44d730cbcd74b777870d473d8e16daccd7ef940b))
* replacing Route with Group, and adding Authentication route ([95ddc70](https://github.com/alvarolucio2007/Scholarly/commit/95ddc70f7bf3fdd62857cf124609541b282bd9bf))

## [3.0.1](https://github.com/alvarolucio2007/Scholarly/compare/v3.0.0...v3.0.1) (2026-10-23)


### Bug Fixes

* replacing writeJSONData with writeJSON in ListX handlers ([26ea1b8](https://github.com/alvarolucio2007/Scholarly/commit/26ea1b883605e0902a206ab3bbbce876cd76803c))
* replacing writeJSONData with writeJSON in ListX handlers ([291109c](https://github.com/alvarolucio2007/Scholarly/commit/291109c7c203b9f09d895b3e399ca4911f4fa9a8))

## [3.0.0](https://github.com/alvarolucio2007/Scholarly/compare/v2.1.1...v3.0.0) (2026-10-23)


### ⚠ BREAKING CHANGES

* Adding proper DTO response to ListAverages and GetAverage
* Adding proper response DTOs and changing all updates to return 204

### minor

* Adding proper DTO response to ListAverages and GetAverage ([ca0be86](https://github.com/alvarolucio2007/Scholarly/commit/ca0be868240620b4d6ca8dffabc350b0f0e7570c))


### Bug Fixes

* Adding proper response DTOs and changing all updates to return 204 ([731e034](https://github.com/alvarolucio2007/Scholarly/commit/731e0347fc45128b3c218f175c54e2542b8fff07))


### Performance Improvements

* Removing unnecessary pointers ([adac3c3](https://github.com/alvarolucio2007/Scholarly/commit/adac3c318ee4adfdfa3ce249b401a79099c84d53))
* Removing unnecessary pointers ([636147c](https://github.com/alvarolucio2007/Scholarly/commit/636147c447704946127226b77beed5765213beef))
* Removing unnecessary pointers and adjusting functions accordingly. ([bb81878](https://github.com/alvarolucio2007/Scholarly/commit/bb81878a8120ada36483b8824953d8315112664d))
* Removing unnecessary pointers from functions ([a5fa47a](https://github.com/alvarolucio2007/Scholarly/commit/a5fa47ac9fe285793334f18641d9026e68f4ae27))
* Removing unnecessary pointers from functions ([aa0f806](https://github.com/alvarolucio2007/Scholarly/commit/aa0f80692a72288b19a83a92c482a6df56021443))
* Removing unnecessary pointers from interfaces ([b629295](https://github.com/alvarolucio2007/Scholarly/commit/b629295c96bb68d59f8c74306548d7e8e600b6c4))
* Removing unnecessary pointers in GetByID and Update ([33383f1](https://github.com/alvarolucio2007/Scholarly/commit/33383f11b038e2dfaa395c0dbfcc1abddf4e5a88))
* Removing unnecessary pointers to satisfy new interface. ([c45e9cb](https://github.com/alvarolucio2007/Scholarly/commit/c45e9cb88b710877c189cc199b0b26ccb908db95))

## [2.1.1](https://github.com/alvarolucio2007/Scholarly/compare/v2.1.0...v2.1.1) (2026-10-20)


### Bug Fixes

* adding nil checking to Create teacher ([09e3862](https://github.com/alvarolucio2007/Scholarly/commit/09e38627eedc0f4d1e74c1c2f09a19847d42a430))

## [2.1.0](https://github.com/alvarolucio2007/Scholarly/compare/v2.0.0...v2.1.0) (2026-10-18)


### Features

* Adding env support ([5e5fcaa](https://github.com/alvarolucio2007/Scholarly/commit/5e5fcaa3ea38d423129ba2a9c986519ea9030311))


### Bug Fixes

* Adding error handling to LoadEnv ([fb7d2e3](https://github.com/alvarolucio2007/Scholarly/commit/fb7d2e3ad25d62fb65cbfdef03f983a8e4a83241))
* Changing the alive conn and idle conn slightly in Postgres.NewConn ([0bef373](https://github.com/alvarolucio2007/Scholarly/commit/0bef373044bd4c038b759e551b72d0980c6ba00f))

## [2.0.0](https://github.com/alvarolucio2007/Scholarly/compare/v1.2.0...v2.0.0) (2026-10-16)


### ⚠ BREAKING CHANGES

* replacing EnrollmentStatus with SituationStatus in ReportCard
* Changing vw_student average return to 'approved', 'recovery' and 'failed'

### Bug Fixes

* Changing vw_student average return to 'approved', 'recovery' and 'failed' ([04fda07](https://github.com/alvarolucio2007/Scholarly/commit/04fda07efb96d66381cca9b155444ac8bb7f509d))
* removing pointer from float64 return in GetAverage ([16b919c](https://github.com/alvarolucio2007/Scholarly/commit/16b919ce99b766ec31921f93b11fed56e3152881))
* replacing EnrollmentStatus with SituationStatus in ReportCard ([0b2d643](https://github.com/alvarolucio2007/Scholarly/commit/0b2d6438fefc24ac3d5d6bfafcb169141806cad0))

## [1.2.0](https://github.com/alvarolucio2007/Scholarly/compare/v1.1.1...v1.2.0) (2026-10-12)


### Features

* Adding package comments to all packages ([a7477e4](https://github.com/alvarolucio2007/Scholarly/commit/a7477e48a1ff76d8814577dc24c29e576aaa7f8d))

## [1.1.1](https://github.com/alvarolucio2007/Scholarly/compare/v1.1.0...v1.1.1) (2026-10-12)


### Bug Fixes

* Adding created_at and updated_at to enrollments migration ([ff6cf73](https://github.com/alvarolucio2007/Scholarly/commit/ff6cf7386fcaaaaf13cf8eca2224a69581399e96))
* Adding updated_at NOW to Update ([2c39203](https://github.com/alvarolucio2007/Scholarly/commit/2c39203a5102ed58903cbe7161a615572b339757))

## [1.1.0](https://github.com/alvarolucio2007/Scholarly/compare/v1.0.0...v1.1.0) (2026-10-10)


### Features

* Add full course cache. ([35cfa82](https://github.com/alvarolucio2007/Scholarly/commit/35cfa82ff6e29b536feb813b62791445926941f8))
* Add full enrollment cache. ([678e3c7](https://github.com/alvarolucio2007/Scholarly/commit/678e3c7e371a1d460d1402d4d023f153bccf1d4e))
* Add full grade cache. ([6e541d9](https://github.com/alvarolucio2007/Scholarly/commit/6e541d9ecd903db7513480cfb232a3dc8517663e))
* Add full student cache ([bac9f82](https://github.com/alvarolucio2007/Scholarly/commit/bac9f82b48d3c1b38ba9619bb95415a82a5fd5c4))
* Add full test cache ([71505c0](https://github.com/alvarolucio2007/Scholarly/commit/71505c0f8de6e7c0c2aa78cd714f36abf3f9559e))
* Adding cache ([7737d29](https://github.com/alvarolucio2007/Scholarly/commit/7737d291583e9c3579eaa7da4721feb0ae0176c0))
* Adding cache ([36aaeb9](https://github.com/alvarolucio2007/Scholarly/commit/36aaeb9191f8bcd835bb34217beaf439d1702ea8))
* Adding cache interfaces. ([59ba137](https://github.com/alvarolucio2007/Scholarly/commit/59ba1376f0644a998e3452909ad22efc3792a55a))
* Adding cache support ([23deeb6](https://github.com/alvarolucio2007/Scholarly/commit/23deeb6516155e625e882687447e783ab3ace88e))
* Adding cache support ([9deb6cc](https://github.com/alvarolucio2007/Scholarly/commit/9deb6cc6256a594624384a6970a1cff75cec8f78))
* Adding cache support ([504b132](https://github.com/alvarolucio2007/Scholarly/commit/504b132a0efa271d9cf07b3cf700b80b3f530d88))
* Adding cache support ([d48a70c](https://github.com/alvarolucio2007/Scholarly/commit/d48a70c5fc0deffdfff5223e21cb6ae633a597ca))
* Adding cache support to student_service ([3099a8d](https://github.com/alvarolucio2007/Scholarly/commit/3099a8d97eabdc432f319f1313e62e5b9e5e484a))
* Adding cache support to teacher_service ([0068de0](https://github.com/alvarolucio2007/Scholarly/commit/0068de03867b46e8e7c4af392ad5813c5109116c))
* Adding CacheQueryTimeout constant ([8c30db7](https://github.com/alvarolucio2007/Scholarly/commit/8c30db7d7d9c2693c1e29ed6d8bef4f7a092bb6c))
* Adding expiry time to user,teacher and student ([12a052e](https://github.com/alvarolucio2007/Scholarly/commit/12a052e36d74247dfcfd1a14c6c93adce1ce6458))
* Adding NewConnn func ([6fc3f5a](https://github.com/alvarolucio2007/Scholarly/commit/6fc3f5a64634216404f6edd457e3bbab7be66945))
* Adding teacher cache ([9910a42](https://github.com/alvarolucio2007/Scholarly/commit/9910a422e1b2b2c28dac4cc11fe4841c1f3cfb28))
* Adding user creation, read and delete in redis ([904a928](https://github.com/alvarolucio2007/Scholarly/commit/904a928bda5b8c93b0a1fc98e21ab476d847cc63))
* Finishing cache in user_service ([f1022b2](https://github.com/alvarolucio2007/Scholarly/commit/f1022b2486efbbe991e946ff54442f55a7c35a54))


### Bug Fixes

* Adding golangci-lint instead of staticcheck (which includes statticcheck) ([e36d8a7](https://github.com/alvarolucio2007/Scholarly/commit/e36d8a75d411370f00968e0b3d30b2723d9fe2e9))
* Adding logger to NewUserService ([d9b645a](https://github.com/alvarolucio2007/Scholarly/commit/d9b645a9c6bf23774e301f3a505bcca251f37fd5))
* changing cacheURL to localhost:6379 ([ca15c0c](https://github.com/alvarolucio2007/Scholarly/commit/ca15c0c4ee1df31318dbfa78ed6fc0bfa49f4e47))
* Changing panic(1) to os.Exit(1) ([a2b8ef5](https://github.com/alvarolucio2007/Scholarly/commit/a2b8ef538dc79d76019105ced389a61f1f7c6a80))
* replacing create user error to create X error message (respective package) ([d3dccc9](https://github.com/alvarolucio2007/Scholarly/commit/d3dccc9a350b5df59203253707fde0d811904dd9))

## 1.0.0 (2026-10-06)


### ⚠ BREAKING CHANGES

* Changing EnrollStudent to receive EnrollStudentDTO body, and adapting func
* Changing function return for CreateTeacher to include *domain.Teacher
* Changing return of CreateStudent and UpdateStudent to have *domain.Student,error, and removing pointer from UserID (int64 now)
* Returning return of CreateStudent and UpdateStudent to have *domain.Student,error, and removing pointer from UserID (int64 now)
* Adding ReportCardFIlter and changing function signatures.

### Features

* Add all enrollment handlers ([81c1530](https://github.com/alvarolucio2007/Scholarly/commit/81c15305ecdc5f6d6b607907e89343bbfc10b927))
* Add all teacherRepo functions ([7e9e271](https://github.com/alvarolucio2007/Scholarly/commit/7e9e27179d445aacb679e0a51130ae3b22eabf11))
* Add all test handlers. ([68d6432](https://github.com/alvarolucio2007/Scholarly/commit/68d6432264f2a244dec5a7d017d3fe68c25b76e7))
* add core user.go file ([f7d0af3](https://github.com/alvarolucio2007/Scholarly/commit/f7d0af338e23b415819667457a4f855a7dc0cc3b))
* Add course domain file ([6fad793](https://github.com/alvarolucio2007/Scholarly/commit/6fad793c873ed5ef7a2e7f545dddde06a98c6158))
* Add enrollment domain file ([0d7cd99](https://github.com/alvarolucio2007/Scholarly/commit/0d7cd9995b7d3d33da7eca57fbac18b9d6a9bd6f))
* Add enrollment number requirement error ([4a1a034](https://github.com/alvarolucio2007/Scholarly/commit/4a1a0343a8db13445c5f906f18bef561fa015d0b))
* Add EnrollStudent function signatre into StudentRepository ([2614222](https://github.com/alvarolucio2007/Scholarly/commit/261422236668eee102db0e42f1c5ea208ec50a41))
* Add grade domain file ([613de16](https://github.com/alvarolucio2007/Scholarly/commit/613de1653127f1cfb200fbfa208d6b0de1b0776d))
* Add invalid email and CPF error cases (returns StatusBadRequest) ([705c052](https://github.com/alvarolucio2007/Scholarly/commit/705c0526ecf06c66bdb6292e9550dd198a7b6ea2))
* Add not null verification to test_date in tests table ([a4b403c](https://github.com/alvarolucio2007/Scholarly/commit/a4b403cf0fdbbfabf727f32f424074959aadfe66))
* Add student domain file ([86d677a](https://github.com/alvarolucio2007/Scholarly/commit/86d677a1f266305bd703bbb60aa166eb2f3a1619))
* Add student domain file ([851a66f](https://github.com/alvarolucio2007/Scholarly/commit/851a66fcaa1aedab1b3be7d4a3861478c63f04bd))
* Add StudentResponseDTO ([4a4ce98](https://github.com/alvarolucio2007/Scholarly/commit/4a4ce98cebe0cfbb13657b9fc3bb33bee7412c94))
* Add swagger docs ; fix: Removing test_date verification from List ([7841e3a](https://github.com/alvarolucio2007/Scholarly/commit/7841e3a75aff8947a4d147c77c58a9edb722994e))
* Add swagger docs to all handlers. ([d477461](https://github.com/alvarolucio2007/Scholarly/commit/d477461682ca01a2fd47a6e27e130e15067e6e51))
* Add swagger godoc to all handlers ([361b36d](https://github.com/alvarolucio2007/Scholarly/commit/361b36d07feb85c4af8d3a9ed7066cdd3b10331e))
* Add swagger godoc to all handlers. ([f1312e5](https://github.com/alvarolucio2007/Scholarly/commit/f1312e5a8c1fdab218cc3d48137f378ee6bc1ec5))
* Add swagger godoc to all handlers. ([4b6abd8](https://github.com/alvarolucio2007/Scholarly/commit/4b6abd88d75a812fb78b33d3de3e4692111ed4ea))
* Add teacher domain file ([8a49385](https://github.com/alvarolucio2007/Scholarly/commit/8a49385c5fb1fe53def62d610ba3317a501d0172))
* Add TeacherResponseDTO ([88a2d76](https://github.com/alvarolucio2007/Scholarly/commit/88a2d76cf10e010eeae69954f5cc2f7c9fdd95d1))
* Add test domain file ([10f2a4e](https://github.com/alvarolucio2007/Scholarly/commit/10f2a4ec289f44de3c81a77511e90d997e6805c1))
* Add TestResponseDTO ([e94e844](https://github.com/alvarolucio2007/Scholarly/commit/e94e844b96f8b7d587d626cfcbb561080e33b40d))
* Add weighed sum func migration files ([c91b6f0](https://github.com/alvarolucio2007/Scholarly/commit/c91b6f0fff8ad4afb72320184705107697b5e575))
* Add weighed_sum sql raw file ([0a7005e](https://github.com/alvarolucio2007/Scholarly/commit/0a7005e4123d3af78d8e3b43ae1556f33fb4d3f6))
* Adding a couple of swagger docs (until ListUsers) ([1350877](https://github.com/alvarolucio2007/Scholarly/commit/1350877e444695bdde853217ba5c048516b76b3a))
* Adding air TOML file ([1cd8e84](https://github.com/alvarolucio2007/Scholarly/commit/1cd8e84027b5c6b6931c5d6b6e079369e70b405b))
* Adding all course handlers ([a8939f5](https://github.com/alvarolucio2007/Scholarly/commit/a8939f59a78bf7b0981a62426e8b75f1c5e19368))
* Adding all courseRepo functions ([f847a50](https://github.com/alvarolucio2007/Scholarly/commit/f847a508252728bce042e28e7852498438751be0))
* Adding all enrollmentsRepo functions ([cecd015](https://github.com/alvarolucio2007/Scholarly/commit/cecd01554c71417e74d438201e4a80fca07da49d))
* Adding all grade handlers. ([c7a5f6d](https://github.com/alvarolucio2007/Scholarly/commit/c7a5f6d57c05e4cd03d64096d3eb5cfe6872f8f4))
* Adding all gradeRepo functions ([8a757ff](https://github.com/alvarolucio2007/Scholarly/commit/8a757ff335daf95000d2341b593274ae86ac7168))
* Adding all routes, to fix. ([4b49bd6](https://github.com/alvarolucio2007/Scholarly/commit/4b49bd65f7ccc89096960be41e035e36375c8827))
* Adding all student handlers ([c8ec721](https://github.com/alvarolucio2007/Scholarly/commit/c8ec7217a9aff2316e69c9bcfdec5ebec329577b))
* Adding all StudentRepository functions ([1ca6138](https://github.com/alvarolucio2007/Scholarly/commit/1ca61384c74e36970b501f24f6eb2cc3960c8161))
* Adding all studentService struct and functions. ([01d992d](https://github.com/alvarolucio2007/Scholarly/commit/01d992dcbf802f97f3753133320aa448d84f3d74))
* Adding all teacher handlers. ([0b24396](https://github.com/alvarolucio2007/Scholarly/commit/0b24396616025976f67065e6f5c443fc27986759))
* Adding all testRepo functions ([c393380](https://github.com/alvarolucio2007/Scholarly/commit/c393380d826d2b2d072567c723cd88b9355180e7))
* Adding all user handlers ([3dac60b](https://github.com/alvarolucio2007/Scholarly/commit/3dac60bbd89f20728f1048d4595b483ff3e06507))
* Adding all UserRepository functions ([ac54814](https://github.com/alvarolucio2007/Scholarly/commit/ac54814f10528c959427c0c4156443a86e82ad1a))
* Adding api container to docker-compose ([bc40a05](https://github.com/alvarolucio2007/Scholarly/commit/bc40a059a67ef16d723a880743d694c4ecd1f48a))
* Adding argon2 password hasher, new, hash and check functions ([fc7d638](https://github.com/alvarolucio2007/Scholarly/commit/fc7d638999da5a8bdd917db3a45c83a868f74448))
* Adding basic config struct ([e565064](https://github.com/alvarolucio2007/Scholarly/commit/e5650649204a88ae90171214b269e9345588cc47))
* Adding basic health check function ([206beac](https://github.com/alvarolucio2007/Scholarly/commit/206beac2683481de915c882aca59c58f27ecd6fb))
* Adding Business rules (WeighedAverage calculation and IsApproved) ([caed20f](https://github.com/alvarolucio2007/Scholarly/commit/caed20f8b9fff8497efe14aeabf4639f5c1bcb69))
* Adding course DTO. ([e135c2d](https://github.com/alvarolucio2007/Scholarly/commit/e135c2d91ed9dd6ea2445eeb69c757d4589bfce9))
* Adding course_dto ([3e5f16d](https://github.com/alvarolucio2007/Scholarly/commit/3e5f16d0c37502a1033b24b580cb3c8baef21b47))
* Adding courseResponseDTO struct ([138d058](https://github.com/alvarolucio2007/Scholarly/commit/138d0583eb99da17595131d9c0ae6ab8eb4930cf))
* Adding courseService functions and struct. ([faa7b22](https://github.com/alvarolucio2007/Scholarly/commit/faa7b2270d8f56ba5e10f2d32b8fe214657116d7))
* Adding CPF verification to Create and Edit functions ([7969d73](https://github.com/alvarolucio2007/Scholarly/commit/7969d737bc331ef35bfd13166a1dcbdffcb9f0c0))
* Adding deleteEnrollment godoc handler, and fixing duplicated route ([43ee50a](https://github.com/alvarolucio2007/Scholarly/commit/43ee50ad9da8c4c8430b762c05f814813be7a562))
* Adding email validation to payload.email ([0af15a5](https://github.com/alvarolucio2007/Scholarly/commit/0af15a599e5d469d8dfd6794c9c8354aa082456c))
* Adding enroll_student procedure migrations ([4b77e00](https://github.com/alvarolucio2007/Scholarly/commit/4b77e00976e6f4c1065c8d6e6f3514296912ef35))
* Adding enrollment DTOs. ([447cecb](https://github.com/alvarolucio2007/Scholarly/commit/447cecb37c8d2a0253c2db5cbf1775dda4a56144))
* Adding enrollment DTOs. ([37a2919](https://github.com/alvarolucio2007/Scholarly/commit/37a2919e7b90d34bca9f31b8554b8782a0f9cc59))
* Adding enrollment service ([052a165](https://github.com/alvarolucio2007/Scholarly/commit/052a1654c74a258f9d27d05b8fa29d486ad5bbff))
* Adding EnrollmentResponseDTO ([c921f00](https://github.com/alvarolucio2007/Scholarly/commit/c921f00e70459ac5c2d2fea903ab748d94b8a4a3))
* Adding EnrollStudent function ([d000ccc](https://github.com/alvarolucio2007/Scholarly/commit/d000cccb9984bbd9455eb2355ba1316a5601ebfd))
* Adding EnrollStudentDTO struc ([0d7cfee](https://github.com/alvarolucio2007/Scholarly/commit/0d7cfee4c30272c0865631a4bae3dadeaf54238e))
* Adding env support for both docker and dev uses. ([9da2441](https://github.com/alvarolucio2007/Scholarly/commit/9da244193dbc3e8b92de31da526d0c23ab18c114))
* Adding env support. ([59a9cd0](https://github.com/alvarolucio2007/Scholarly/commit/59a9cd02cdc578e1da9790394870bc442bcff006))
* Adding errors invalid email and cpf ([25ec534](https://github.com/alvarolucio2007/Scholarly/commit/25ec534a9e0dd74462bd5a7e4a984dcdef63abc1))
* Adding full List function ([41d5fc6](https://github.com/alvarolucio2007/Scholarly/commit/41d5fc671cba25c90b839a921b7b24f38f1bc83f))
* Adding GetAverage and fixing ListAverages godoc [@router](https://github.com/router) ([2a62882](https://github.com/alvarolucio2007/Scholarly/commit/2a628820e869e772a875c96343e45d01470e5983))
* Adding GetAverage function ([2527c11](https://github.com/alvarolucio2007/Scholarly/commit/2527c114356f94fbd205ad59992fa5953389ec11))
* Adding GetAverage router and fixing ListAverages router ([785ca31](https://github.com/alvarolucio2007/Scholarly/commit/785ca3191b4a039b0e16236d2a855e95eb983a31))
* Adding GetAverage to GradeRepository and renaming ListReportCard to ListAverages ([8f69cd9](https://github.com/alvarolucio2007/Scholarly/commit/8f69cd9156a42f529c64320b48a48837fc747a17))
* Adding GetAverage, renaming ListReportCard to ListAverages and adding proper error translation to every error. ([6e105a0](https://github.com/alvarolucio2007/Scholarly/commit/6e105a0efba1ea81017e4a14f7fb775a278f3b90))
* Adding GetWeighedAverage function ([47de855](https://github.com/alvarolucio2007/Scholarly/commit/47de8556bd62c6e7acd0bacb529808948c2f2cdb))
* Adding golang migrations directly into the docker-compose ifle ([ab6f9ed](https://github.com/alvarolucio2007/Scholarly/commit/ab6f9ede213f42e44ae74386d494c4734320e005))
* Adding grade DTOs. ([79f3a55](https://github.com/alvarolucio2007/Scholarly/commit/79f3a558d13749649f3ab60f16382b89b01a65ff))
* Adding grade handler ([9ba11b4](https://github.com/alvarolucio2007/Scholarly/commit/9ba11b4a375728e640e3688aeb5b74c4918a0208))
* Adding GradeResponseDTO ([1c965d5](https://github.com/alvarolucio2007/Scholarly/commit/1c965d5086bcbab81a96b8e0c32a2e0a5fc90b30))
* Adding handler struct (with all services) and NewHandler func ([2b81e18](https://github.com/alvarolucio2007/Scholarly/commit/2b81e188e39911910704a16920d4ba0c8ac2774f))
* adding JSON-specific functions (writeJSON,writeJSONError,writeJSONData,respondDomainError) ([e550894](https://github.com/alvarolucio2007/Scholarly/commit/e5508947682a192c3661aac3aed46539a4f6e591))
* Adding logger package, New function which returns *slog.Logger ([d852e5d](https://github.com/alvarolucio2007/Scholarly/commit/d852e5d80cb52657592d10a543a57713faca6054))
* Adding multi-stage dockerfile ([7308e3e](https://github.com/alvarolucio2007/Scholarly/commit/7308e3e920ab28fda98ed2aebf3ded8dae7bf005))
* Adding NewGradeRepo function ([70b151d](https://github.com/alvarolucio2007/Scholarly/commit/70b151d8cbd530369044cc9ca910e63326dc8ca8))
* Adding nilIfEmpty function ([a051cd8](https://github.com/alvarolucio2007/Scholarly/commit/a051cd80bd4ef840135ff8b17d3c9e787bc211c0))
* Adding nilIfZeroInt and nilIfZeroFloat funcs ([82d5958](https://github.com/alvarolucio2007/Scholarly/commit/82d595841e48d2368b15da94ad287341ace74a7d))
* Adding PasswordHasher ([5c42826](https://github.com/alvarolucio2007/Scholarly/commit/5c428268d64551d646b9d50fdedbbaa4accd44c5))
* Adding postgres and postgres repo connections ([a5f1803](https://github.com/alvarolucio2007/Scholarly/commit/a5f1803b1f6720af1a20e56572751481e2626016))
* Adding postgres connection function. ([a35a766](https://github.com/alvarolucio2007/Scholarly/commit/a35a76618390ea54d30c4a822509ac090fb3583c))
* Adding proper structured logging ([3ace38b](https://github.com/alvarolucio2007/Scholarly/commit/3ace38bc5258aa9c5633109b72aa93275b68277a))
* Adding release please action workflow ([e90d867](https://github.com/alvarolucio2007/Scholarly/commit/e90d867e315d1a4085e34129da175c54b6227aff))
* Adding release please action workflow ([1fc1151](https://github.com/alvarolucio2007/Scholarly/commit/1fc11515558b6ef9b8aeeced432671218c8f7572))
* Adding ReportCardFIlter and changing function signatures. ([38e4d5a](https://github.com/alvarolucio2007/Scholarly/commit/38e4d5a90039d0c69490b2ee21edcfde4699753f))
* Adding ReportCardRepository interface ([583b27b](https://github.com/alvarolucio2007/Scholarly/commit/583b27b1b867ef8bb924db40fc5810fd93709cc8))
* adding ReportCardResponseDTO and updating updateGradeDTO to exclude pointers and fixing JSON tags ([aee168c](https://github.com/alvarolucio2007/Scholarly/commit/aee168c52ac13897ab46f8cefdbcd3f2dd53fca8))
* Adding repositories interfaces ([cc2ffd3](https://github.com/alvarolucio2007/Scholarly/commit/cc2ffd32b13f31474ca9d82bbd0aa42f74cba888))
* Adding repositories interfaces ([4352e25](https://github.com/alvarolucio2007/Scholarly/commit/4352e25884484b319b620dbdb32ea4adc507d6a3))
* Adding run command ([b2b01ce](https://github.com/alvarolucio2007/Scholarly/commit/b2b01ce8996b8ba988572b1cf51fba4634661f99))
* Adding slog to the handler ([16c4db2](https://github.com/alvarolucio2007/Scholarly/commit/16c4db249358f0b2356599c098fbbb04e31c661c))
* Adding sp_enroll_student procedure ([2fc95b7](https://github.com/alvarolucio2007/Scholarly/commit/2fc95b7f0eb566981211154f81b4ecc9f38f532f))
* Adding student report card view ([2230616](https://github.com/alvarolucio2007/Scholarly/commit/2230616dfbbcbb0b0a027246a42ddd0ba9c0dc84))
* Adding student report card view migration ([a768552](https://github.com/alvarolucio2007/Scholarly/commit/a768552d80a7b250a19033ce9560532a3730448d))
* Adding swagger gen-docs command ([aa758eb](https://github.com/alvarolucio2007/Scholarly/commit/aa758ebc834e5f1209011ff40db96237001e6feb))
* Adding teacher DTO ([a25ac9d](https://github.com/alvarolucio2007/Scholarly/commit/a25ac9d8097b40a2982d71e3753979a201e16d8a))
* Adding teacherService functions and struct ([60bb60a](https://github.com/alvarolucio2007/Scholarly/commit/60bb60a1a6bbd79b666544099008d6b0a1fd873c))
* Adding template files for better structuting ([99fe816](https://github.com/alvarolucio2007/Scholarly/commit/99fe81692eb4f522ed788fc65cfbedb70cc589b0))
* Adding test DTOs. ([6fdd6fc](https://github.com/alvarolucio2007/Scholarly/commit/6fdd6fc8c899d63346a32add4202d7eedde222cf))
* Adding test service ([ceceb0c](https://github.com/alvarolucio2007/Scholarly/commit/ceceb0cd090a4cbbf538a24ac7476c0569890b70))
* Adding the remainder of main ([f329cfe](https://github.com/alvarolucio2007/Scholarly/commit/f329cfe93cae0561ee8c18c6fe7978b5096447d7))
* Adding the remainder of slogging ([7facd66](https://github.com/alvarolucio2007/Scholarly/commit/7facd66522199dee6275bad82ba846dd0eb377e8))
* Adding the remainder of the swagger docs ([9b4e88d](https://github.com/alvarolucio2007/Scholarly/commit/9b4e88d3e200cd21ae0e5f4ae83c883598cadb0d))
* Adding user DTOs ([b682ad5](https://github.com/alvarolucio2007/Scholarly/commit/b682ad55e7a1f7cb71349ae9a94a73d4e6ca34a9))
* Adding user_service ([54f731a](https://github.com/alvarolucio2007/Scholarly/commit/54f731a15a6cc1c374bfd7bfa014bf66d4e346ee))
* Adding user_service service file ([958b563](https://github.com/alvarolucio2007/Scholarly/commit/958b563d34d93d1d913dcb4f1c3f331145bc4af6))
* adding UserResponseDTO to here ([f548ca7](https://github.com/alvarolucio2007/Scholarly/commit/f548ca71cdb4a5daff8b7a16b8ae93e3ad989459))
* adding validateEmail and proper CPF validation func ([8475c69](https://github.com/alvarolucio2007/Scholarly/commit/8475c69a7124d64fc3373fcaf5205b1dc72841ca))
* Adding validator ([ec585a8](https://github.com/alvarolucio2007/Scholarly/commit/ec585a8c817f2bbd9a14bc7ed86c7c0f283f9b31))
* Adding very basic router, to be implemented with all handlers. ([f8f97e5](https://github.com/alvarolucio2007/Scholarly/commit/f8f97e5f30570f883c99993b5991a8c6ba10821e))
* Adding View Report card postgres skeleton. ([8d79ab3](https://github.com/alvarolucio2007/Scholarly/commit/8d79ab3cf18aaab69e46d1fdb8ae69e41744ecbe))
* Addnig CreateStudent and UpdateStudent DTOs ([b75c47f](https://github.com/alvarolucio2007/Scholarly/commit/b75c47f4c44aaf94e4f622cbecf145038f923d0c))
* Changing mount and run functions to public, so they can be used in main ([d43e0dc](https://github.com/alvarolucio2007/Scholarly/commit/d43e0dc65ce016cd2d29333fccf74870df8a9b11))
* Implementing hello world Go prototype. ([1d6cd09](https://github.com/alvarolucio2007/Scholarly/commit/1d6cd09efbdbdd7f71f9a863127ea75f94a35ee1))
* Removing status from CreateEnrollmentDTO ([c52c54a](https://github.com/alvarolucio2007/Scholarly/commit/c52c54a14be45253e903574d361ce24185a54b1c))


### Bug Fixes

* Adapting enrollStudent route to receive metadata in body ([cd60e30](https://github.com/alvarolucio2007/Scholarly/commit/cd60e30a469b2c5176183d0c41ea464e46b58922))
* Added proper error translation to funcitons ([e119b54](https://github.com/alvarolucio2007/Scholarly/commit/e119b54e151fee48fbaf92e59b7aa522a7d5736b))
* Adding CourseID to UpdateTest and respective payload ([d1da7a2](https://github.com/alvarolucio2007/Scholarly/commit/d1da7a2bc7a47b495afc2a1cc22154031264da5d))
* adding error to InternalServerError ([527f7e3](https://github.com/alvarolucio2007/Scholarly/commit/527f7e32738d9ea7681e3c9efab106e70a849581))
* Adding error translation to several other functions ([a8cd029](https://github.com/alvarolucio2007/Scholarly/commit/a8cd029f391166dc68c45b05b0f751b88862263a))
* Adding final $$ ([e4ab241](https://github.com/alvarolucio2007/Scholarly/commit/e4ab241733590a8cea707b3d346ba16873d7e946))
* Adding logger to NewUserRepo ([b3c83c7](https://github.com/alvarolucio2007/Scholarly/commit/b3c83c75a53591aba60d665f6f92b45cdd0627b5))
* adding logger to userRepo ([f4a4c18](https://github.com/alvarolucio2007/Scholarly/commit/f4a4c1877aeb61f2de4bd6abaf1d981c533a90c9))
* Adding proper route usage in httpSwagger.Handler, and fixing a couple of routes ([662905c](https://github.com/alvarolucio2007/Scholarly/commit/662905cf12a571b0f75693631cf3680c06f14a39))
* Adding return to writeJSONError in GetEnrollmentByID func ([3cf0991](https://github.com/alvarolucio2007/Scholarly/commit/3cf0991414d5bf125ffbaa5f38cf0bf7b109fb73))
* Adding updated_at to Update query ([f3919f4](https://github.com/alvarolucio2007/Scholarly/commit/f3919f4da9caf3dd0d6ba17e60b008b12aa3db4d))
* adding UserResponse struct and making CreateUser return *UserResponse ([d6185c4](https://github.com/alvarolucio2007/Scholarly/commit/d6185c41cb6eba0f9f24ed86645280268829d37d))
* Adidng proper error handling in writeJSONError, testDate parsing in CreateTest, and adding testDate parsing in UpdateTest ([93b6d10](https://github.com/alvarolucio2007/Scholarly/commit/93b6d10dc2f71c16a493ff8e2c4f69589f3deace))
* Changing EnrollStudent to receive EnrollStudentDTO body, and adapting func ([bc17ca3](https://github.com/alvarolucio2007/Scholarly/commit/bc17ca3cce039047d2a923ebdd3757c22004abe6))
* changing function GetWeighedAverage to GetWeightedAverage and fixing procedure call in EnrollStudent ([9814682](https://github.com/alvarolucio2007/Scholarly/commit/981468288c23c2da120ae512c87e63d4441fce4d))
* Changing function return for CreateTeacher to include *domain.Teacher ([50e84da](https://github.com/alvarolucio2007/Scholarly/commit/50e84da238f387d91c515fd5d38749fd6b82b43a))
* Changing return of CreateStudent and UpdateStudent to have *domain.Student,error, and removing pointer from UserID (int64 now) ([07cd800](https://github.com/alvarolucio2007/Scholarly/commit/07cd800829a03468177a3e5090c75fb7444e1249))
* changing route swagger doc in GetUserByID to /users/ ([4693bb3](https://github.com/alvarolucio2007/Scholarly/commit/4693bb3eb3f59e970249a5ea2356f9d56c358f46))
* changing StatusOK to StatusCreated in EnrollStudent ([bd240b2](https://github.com/alvarolucio2007/Scholarly/commit/bd240b2a9eb9e93acc73e141023718365f34adb9))
* Changing tags to correct ones ([0575c59](https://github.com/alvarolucio2007/Scholarly/commit/0575c590c8d2028c20c095e0d968db22345194e4))
* Changing tags to correct ones. ([42135cc](https://github.com/alvarolucio2007/Scholarly/commit/42135cc0dde8c694fa2575f5b6b64e24059222a8))
* Changing Test to domain.Test in GetTestByID godoc ([481a451](https://github.com/alvarolucio2007/Scholarly/commit/481a45147a12107340bba5356662c7f09cfbab11))
* changing TestDate from time.Time to string in CreateTestDTO and UpdateTestDTO, and updating it ([f34f156](https://github.com/alvarolucio2007/Scholarly/commit/f34f156a7afb5308b00d687104d90fd82829c6b6))
* Changing UserDTO responses to TeacherDTO responses (also another minor fixes) in GODOC ([6676ec3](https://github.com/alvarolucio2007/Scholarly/commit/6676ec33028b52fcdbf5c1e3161a0a79afdc27a7))
* Changing valid CPF with punctuation to false in test (it needs to receive without punctuation) ([6451f9e](https://github.com/alvarolucio2007/Scholarly/commit/6451f9ee6f3f0dab1ba010e009ccfb27e5a2ef50))
* Ensuring proper spacing between godoc and respective functions ([e43e596](https://github.com/alvarolucio2007/Scholarly/commit/e43e5964daa2997377211b83ae6f733d947ef44f))
* Fix UpdateTest logic ([aabb843](https://github.com/alvarolucio2007/Scholarly/commit/aabb8431107ed65ec624b5b095935161639ffd09))
* Fixing  mount by adding CORS and proper swagger handler. ([27ec576](https://github.com/alvarolucio2007/Scholarly/commit/27ec57671b836b2e9c86df626a16fbb9de996fb5))
* Fixing create and getByEnrollment functions. ([b5073ee](https://github.com/alvarolucio2007/Scholarly/commit/b5073ee11230a0c5cf2605437f47e78e3cd0345c))
* Fixing duplicated route and changing variable name to more Go-esque. ([c0bdcf4](https://github.com/alvarolucio2007/Scholarly/commit/c0bdcf4c862af3dbe582e8f2da59b3f35553fc8e))
* Fixing duplicated routes in godoc ([2cfce6b](https://github.com/alvarolucio2007/Scholarly/commit/2cfce6b197f43df7f4bcdbf182a0649e8ff16549))
* fixing error not checking in writeJSONData ([d3c99a2](https://github.com/alvarolucio2007/Scholarly/commit/d3c99a2481ff85d84b8ff6a0ed6485aa0572df7d))
* Fixing GetStudentByEnrollment route ([2d5603c](https://github.com/alvarolucio2007/Scholarly/commit/2d5603cec40a8a09e868fe21316af9eaeafdb9af))
* Fixing GetStudentBYID handler and godoc ([6cfa4b4](https://github.com/alvarolucio2007/Scholarly/commit/6cfa4b4d558fe39f255d7d8f399130179cbc1331))
* Fixing ListReportCard godoc handler and removing & from final writeJSONData in same func ([a00e1f6](https://github.com/alvarolucio2007/Scholarly/commit/a00e1f6795f9cac0111f9dcadc8238ec478e712b))
* Fixing route in GetTestByID ([79bc6fd](https://github.com/alvarolucio2007/Scholarly/commit/79bc6fd49c73d8be55e35ebae0cf2e5c0af1d8f1))
* Fixing route in ListReportCard ([23d1236](https://github.com/alvarolucio2007/Scholarly/commit/23d12362573af791c7b5ab88d021bea04ef49553))
* Fixing the project skeleton. ([8eca5b1](https://github.com/alvarolucio2007/Scholarly/commit/8eca5b1db13cffce66a9dc7708bbe462bddbc3c5))
* Fixing UpdateTestDTO struct logic (removing pointers,fixing validates) ([db4d1e5](https://github.com/alvarolucio2007/Scholarly/commit/db4d1e527c7cab60c9461327fad8534e6a3ed3fa))
* making cors allowed origins into *, and adding addr string as parameter in Run ([2d85f61](https://github.com/alvarolucio2007/Scholarly/commit/2d85f617e550009c7fcec9bbe204b946b28c4f7b))
* Removing nonsensical quotation marks from integrers in Update (student_id and course_id) ([43b3f01](https://github.com/alvarolucio2007/Scholarly/commit/43b3f0163a15d11b23c96abd7e03df4473da2a5a))
* Removing pointers from name-passowrd in UpdateUserDTO, and making example values "" ([a9865ba](https://github.com/alvarolucio2007/Scholarly/commit/a9865ba7ed6ddc055decdf866fe51e861c746390))
* Removing trailing / in ListAverages godoc ([65107eb](https://github.com/alvarolucio2007/Scholarly/commit/65107eb4f7eeff203e7168c0c1e4bec87ce61e9d))
* Removing unnecessary non Business-logic (Only CPF verification remains) ([b47911c](https://github.com/alvarolucio2007/Scholarly/commit/b47911c6ff02c0316fdf25707251e1d3bfab4559))
* Removing unused addIf function. ([924af9e](https://github.com/alvarolucio2007/Scholarly/commit/924af9e1be033fd018db40b6bbdab73cdfea2493))
* replace users with test ([89a40d2](https://github.com/alvarolucio2007/Scholarly/commit/89a40d26143fc877c48d18bddcc8463fa89bf6c1))
* replacing domain.Grade to GradeResponseDTO ([98cafb0](https://github.com/alvarolucio2007/Scholarly/commit/98cafb0f6cbefefc518dd0611c08aaf832234435))
* Replacing domain.Student with StudentResponseDTO ([94aaa43](https://github.com/alvarolucio2007/Scholarly/commit/94aaa4311b0091fa8df99dfa161f22f3c5a7e017))
* replacing domain.Teacher with TeacherResponseDTO ([6be7938](https://github.com/alvarolucio2007/Scholarly/commit/6be79389fddb179b00648e6e78b06adc5b26eb97))
* Replacing domain.Test with TestResponseDTO ([6299114](https://github.com/alvarolucio2007/Scholarly/commit/6299114e8a44cfffeae6a4d25a0ef18666b74ca9))
* replacing panic(err) with panic(1) ([77769dd](https://github.com/alvarolucio2007/Scholarly/commit/77769dd6de70d5cf09ba65d8c8fcd579b5ef5afa))
* Returning return of CreateStudent and UpdateStudent to have *domain.Student,error, and removing pointer from UserID (int64 now) ([c7d21a1](https://github.com/alvarolucio2007/Scholarly/commit/c7d21a17adeb07a2ca5c0b6bee40b1d5e4326cd1))
* Reverting CreateUser to *domain.User,error return ([0aeee5e](https://github.com/alvarolucio2007/Scholarly/commit/0aeee5ec535f5941835369073fd7ab8b0c2b59fe))
* Switching domain.Course to CourseResponseDTO ([5ec2268](https://github.com/alvarolucio2007/Scholarly/commit/5ec22686430ab457f25185740682bc1799a08cc5))
