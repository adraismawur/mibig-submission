from wtforms import (
    BooleanField,
    Form,
    FormField,
    HiddenField,
    IntegerField,
    PasswordField,
    StringField,
    SubmitField,
    validators,
)
from wtforms.widgets import HiddenInput

from submission.models.users import UserInfo

class UserInfoForm(Form):
    db_id = IntegerField(widget=HiddenInput())
    db_user_id = IntegerField(widget=HiddenInput())

    alias = StringField(widget=HiddenInput(), default=UserInfo.generate_alias())

    name = StringField(
        "Name",
        validators=[
            validators.InputRequired(message="Please provide a name")
        ],
    )

    call_name = StringField(
        "What should we call you?",
        validators=[
            validators.InputRequired(message="Please provide a way to address you")
        ],
    )
    orc_id = StringField("ORCID")
    organisation_1 = StringField(
        "First affiliation",
        validators=[
            validators.InputRequired(message="Please provide an affiliation")
        ],
    )
    organisation_2 = StringField("Second affiliation (optional)")
    organisation_3 = StringField("Third affiliation (optional)")
    public = BooleanField(widget=HiddenInput())


class UserDetailsEditForm(Form):
    db_id = IntegerField(widget=HiddenInput())
    # anonymous = BooleanField("Anonymize me", description="Set yourself as an anonymous contributor. When checked this will replace your unique user ID with an anonymous one.",)

    email = StringField(widget=HiddenInput())

    password = PasswordField("Change password", description="Leave blank to leave un-changed.")
    password_confirm = PasswordField("Confirm password", validators=[validators.EqualTo("password")])

    info = FormField(UserInfoForm)
    
    submit = SubmitField("Update user details")
