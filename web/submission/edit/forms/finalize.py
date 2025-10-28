from wtforms import (
    Form,
    SelectField,
    BooleanField,
    StringField,
    SubmitField,
    validators,
)
from submission.utils.custom_widgets import SelectDefault, SubmitIndicator


class FinalizeForm(Form):

    completeness = SelectField(
        "Completeness *",
        choices=["complete", "incomplete", "unknown"],
        description="Are all genes needed for production of compounds present in the specified locus/loci?",
        widget=SelectDefault(),
        validate_choice=False,
        validators=[validators.InputRequired()],
    )
    embargo = BooleanField(
        description="Please embargo my gene cluster information, pending publication of the results. "
        "For newly characterized gene clusters only. Please notify us upon publication so that the embargo can be lifted."
    )
    comments = StringField("Additional comments")

    submit = SubmitField("Submit", widget=SubmitIndicator())
