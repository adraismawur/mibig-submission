from wtforms import (
    Form,
    PasswordField,
    EmailField,
    BooleanField,
    StringField,
    SubmitField,
    validators,
)
from wtforms.widgets import HiddenInput


class LoginForm(Form):
    username = EmailField("Username", render_kw={"placeholder": "alice@example.edu"})
    password = PasswordField("Password")
    remember = BooleanField("Remember me")
    submit = SubmitField("Login")


class UserEmailForm(Form):
    email = EmailField(
        "Email",
        validators=[
            validators.InputRequired(message="Please provide an email address")
        ],
        render_kw={"placeholder": "alice@example.edu"},
    )
    submit = SubmitField("Send email")


class PasswordResetForm(Form):
    email = StringField(render_kw={"disabled": ""})
    new_password = PasswordField("New password")
    challenge = StringField(widget=HiddenInput())
    confirm = PasswordField(
        "Confirm password",
        validators=[validators.EqualTo("new_password", "Password mismatch!")],
    )
    submit = SubmitField("Save password")
