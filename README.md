### Тестирование

Решил добавить JWT-сервис на одном уровне с REPO-сервисом, чтобы была возможность тестировать сервис аутентификации.

Также, было много сомнений по поводу отсутсвия локиги в бизнес-слое, но по-другому сделать никак не получится для такого контракта интерфейса, соотвественно тестировать невозможно большую часть бизнес-логики. Были моменты, когда поиск по параметру в БД (например, поиск юзера по имени) можно вынести в цикл в бизнес-логике, а не оставлять это в SQL. Тогда и сам слой БД будет выглядеть более "строгим" и добавится логика в бизнес-слое, но опять тестированию это не поможет

### Хэндлеры

Я не думаю, что хэндлеры выглядят весьма большими, проверка в каждом из них токена - ведь не катастрофа, да?

### Репозиторий

Старался выделять 