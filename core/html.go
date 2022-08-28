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

	tgAuth = `<div class="form_auth_block">
  <div class="form_auth_block_content">
    <p class="form_auth_block_head_text">Авторизуйтесь в телеграм</p>
    <form class="form_auth_style" action="/login" method="post">
      <label>Номер телефона</label>
      <input type="phone" name="phone" placeholder="89009009090" required >
            <button class="form_auth_button" type="submit" name="form_auth_submit">Продолжить</button>
    </form>
  </div>
</div>`

	tgAuthCode = `<div class="form_auth_block">
  <div class="form_auth_block_content">
    <p class="form_auth_block_head_text">Введите проверочный код</p>
    <form class="form_auth_style" action="/login/submit" method="post">
      <label>Код</label>
      <input type="numbers" name="code" placeholder="12345" required >
            <button class="form_auth_button" type="submit" name="form_auth_submit">Продолжить</button>
      <input value="%s" hidden name="phone" placeholder="12345" readonly >
      <input value="%s" hidden name="hash" placeholder="12345" readonly >
    </form>
  </div>
</div>`

	tgPasswd = `<div class="form_auth_block">
  <div class="form_auth_block_content">
    <p class="form_auth_block_head_text">Двухфакторная авторизация</p>
    <form class="form_auth_style" action="password" method="post">
      <label>Пароль двухфакторной авторизации</label>
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
  </div>
</div>`
)
