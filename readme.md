**Home-counter** - это сервис для подсчета платы за коммунальные услуги. 
Вам необходимо залогиниться, указать текущие показания счетчиком и тарифы.
Дальше каждый месяц вы указываете новые показания счетчиков и получаете информацию о том, сколько нужно заплатить
за коммуналку в этом месяце.

*Фронтенд пока не готов. Ждет Дениса)*

##API

- /auth/login (GET)
- /auth/logout (GET)

#### В следующих ручках нужна кука, которую проставит сервис после логина
- /user (GET) вернет инфу о пользаке
- /user/tariffs (POST) отправляем в json-e текущие тарифы (electricity, hot_water,
  cold_water, out_water, internet)
- /user/meters (POST) отправляем в json-e показания счетчиков (electricity, hot_water,
  cold_water)
- /user/count-meters (GET) узнаем, сколько должны за коммуналку (параметры electricity, hot_water,
  cold_water)

*Подробная дока по работе с миграциями: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate*
Накатить новые миграции:
`$ migrate -source file://path/to/migrations -database postgres://localhost:5432/database up`
 Откатить последнюю миграцию
`$ migrate -source file://path/to/migrations -database postgres://localhost:5432/database down 1`