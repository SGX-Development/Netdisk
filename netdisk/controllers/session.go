package controllers

//这个结构体有问题，后面修
type UserStatus struct {
	userName string
	userid   int
	islogin  bool
}

func (c *MainController) Islogin() bool {
	status := c.GetSession("status")
	if status == nil {
		c.Redirect("/login", 302)
	}
	return !(status == nil || (status != nil && !status.(UserStatus).islogin))
}

func (c *IntroController) Islogin() bool {
	status := c.GetSession("status")
	if status == nil {
		c.Redirect("/login", 302)
	}
	return !(status == nil || (status != nil && !status.(UserStatus).islogin))
}

func (c *RecordController) Islogin() bool {
	status := c.GetSession("status")
	if status == nil {
		c.Redirect("/login", 302)
	}
	return !(status == nil || (status != nil && !status.(UserStatus).islogin))
}

func (c *AboutController) Islogin() bool {
	status := c.GetSession("status")
	if status == nil {
		c.Redirect("/login", 302)
	}
	return !(status == nil || (status != nil && !status.(UserStatus).islogin))
}

func (c *BinController) Islogin() bool {
	status := c.GetSession("status")
	if status == nil {
		c.Redirect("/login", 302)
	}
	return !(status == nil || (status != nil && !status.(UserStatus).islogin))
}

func (c *MainController) UserName() string {
	status := c.GetSession("status")
	if status == nil {
		c.Redirect("/login", 302)
	}
	return status.(UserStatus).userName
}

func (c *RecordController) UserName() string {
	status := c.GetSession("status")
	if status == nil {
		c.Redirect("/login", 302)
	}
	return status.(UserStatus).userName
}

func (c *IntroController) UserName() string {
	status := c.GetSession("status")
	if status == nil {
		c.Redirect("/login", 302)
	}
	return status.(UserStatus).userName
}

func (c *AboutController) UserName() string {
	status := c.GetSession("status")
	if status == nil {
		c.Redirect("/login", 302)
	}
	return status.(UserStatus).userName
}

func (c *BinController) UserName() string {
	status := c.GetSession("status")
	if status == nil {
		c.Redirect("/login", 302)
	}
	return status.(UserStatus).userName
}

func (c *RegController) Islogin() bool {
	status := c.GetSession("status")
	if status == nil {
		c.Redirect("/login", 302)
	}
	return !(status == nil || (status != nil && !status.(UserStatus).islogin))
}
