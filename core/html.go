package core

const (
	login = `<div class="form_auth_block">
  <div class="form_auth_block_content">
    <p class="form_auth_block_head_text">Авторизация</p>
    <form class="form_auth_style" action="#" method="post">
      <label>Логин</label>
      <input type="text" name="login" placeholder="Введите логин" required >
      <label>Пароль</label>
      <input type="password" name="password" placeholder="Введите пароль" required >
	  <button class="form_auth_button" type="submit" name="form_auth_submit">Продолжить</button>
    </form>
  </div>
</div>`

	cities = `<div class="form_auth_block">
  <div class="form_auth_block_content">
    <p class="form_auth_block_head_text">Выберите город</p>
		<form action="/volgograd" method="get">
			<input type="submit" value="Волгоград" />
		</form>
		<form action="/moscow" method="get">
			<input type="submit" value="Москва" />
		</form>
		<form action="/krasnodar" method="get">
			<input type="submit" value="Краснодар" />
		</form>
		<form action="/logout" method="get">
			<input type="submit" value="Выйти" />
		</form>
  </div>
</div>`

	fmtReportBody = `<div class="form_auth_block">%s</div>`

	fmtMedia = `        <div class="image">
            <img src="%s" />
        </div>`

	fmtText = `        <div class="info">
			<p>%s</p>
        </div>`
	fmtReport = `<article class="grid-item">
			<h3>Чат: %s, От: @%s, Время: %s</h3>
				%s
    </article><hr/>`
)
