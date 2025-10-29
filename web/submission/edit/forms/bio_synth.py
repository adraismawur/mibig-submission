from flask import current_app, session
import requests
from wtforms import (
    Form,
    StringField,
    IntegerField,
    FieldList,
    SelectField,
    SelectMultipleField,
    FormField,
    BooleanField,
    ValidationError,
    validators,
    SubmitField,
)
from submission.utils.custom_fields import TagListField
from submission.utils.custom_forms import location_form_factory, EvidenceForm
from submission.utils.custom_widgets import (
    FieldListAddBtn,
    TextInputIndicator,
    SubmitIndicator,
    SelectDefault,
    ProductInputSearch,
)
from submission.utils.custom_validators import (
    ValidateSingleInput,
    validate_genbank,
    validate_loci,
)


class BiosynthClassForm(Form):
    class_ = StringField(id="class", name="class")
    subclass = StringField()
    cyclases = FieldList(
        StringField(),
        widget=FieldListAddBtn(
            label="Add cyclase",
        ),
    )


class BiosynthModuleForm(Form):
    type = StringField()
    genes = FieldList(StringField(), widget=FieldListAddBtn(label="Add gene"))
    active = BooleanField()


class BiosynthForm(Form):
    classes = FieldList(
        FormField(BiosynthClassForm),
        min_entries=1,
        description="List of classes in this entry",
        widget=FieldListAddBtn(
            label="Add class",
        ),
    )

    modules = FieldList(
        FormField(BiosynthModuleForm), widget=FieldListAddBtn(label="Add module")
    )


class BioSynthForm(Form):

    biosynthesis = FormField(BiosynthForm)

    submit = SubmitField("Continue to ...", widget=SubmitIndicator())
