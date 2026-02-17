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

    acknowledge = BooleanField(
        "Confirm",
        description="By checking this you indicate that you have, to the best of your abilities, confirmed "
        "that the information in this submission is correct",
        validators=[validators.InputRequired()],
    )

    submit = SubmitField("Submit for review", widget=SubmitIndicator())
