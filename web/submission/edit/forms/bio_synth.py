from flask import current_app, session
import requests
from wtforms import (
    Form,
    HiddenField,
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
from wtforms.widgets import HiddenInput
from submission.edit.forms.biosynthesis import OperonForm
from submission.edit.forms.biosynthesis_paths import PathForm
from submission.utils.custom_fields import TagListField
from submission.utils.custom_forms import location_form_factory, LociEvidenceForm
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

biosynthetic_classes = [
    "NRPS",
    "terpene",
    "PKS",
    "ribosomal",
    "saccharide",
    "other",
]

biosynthetic_sub_classes = [
    "aminocoumarin",
    "aminoglycoside",
    "butyrolactone",
    "cyclitol",
    "Diterpene",
    "ectoine",
    "exopolysaccharide",
    "fatty acid",
    "flavin",
    "Hemiterpene",
    "hybrid/tailoring",
    "Iterative type I",
    "lipopolysaccharide",
    "Modular type I",
    "Monoterpene",
    "non-nrp beta-lactam",
    "non-nrp siderophore",
    "nucleoside",
    "oligosaccharide",
    "other",
    "pbde",
    "phenazine",
    "phosphonate",
    "RiPP",
    "Sesquiterpene",
    "Sesterterpene",
    "shikimate-derived",
    "Trans-AT type I",
    "Triterpene",
    "trna-derived",
    "Type I",
    "Type II aromatic",
    "Type II arylpolyene",
    "Type II highly reducing",
    "Type II",
    "Type III",
    "Type IV",
    "Type V",
    "Type VI",
    "Unknown",
    "unmodified",
]


class BiosynthClassForm(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0)
    db_biosynth_id = IntegerField(widget=HiddenInput(), default=0)
    class_ = SelectField(id="class", name="class", choices=biosynthetic_classes)
    subclass = SelectField(choices=biosynthetic_sub_classes)
    cyclases = FieldList(
        StringField(),
        widget=FieldListAddBtn(
            label="Add cyclase",
        ),
    )


module_types = [
    "pks-modular",
    "nrps-type1",
    "pks-trans-at",
]


class ModuleLocationForm(Form):
    from_ = IntegerField()
    to = IntegerField()


class BiosynthModuleForm(Form):
    db_id = IntegerField(widget=HiddenInput())
    db_biosynth_id = IntegerField(widget=HiddenInput(), default=0)
    genes = FieldList(
        StringField(default="Gene ID"),
        widget=FieldListAddBtn(label="Add gene"),
    )
    name = StringField()
    type = SelectField(choices=module_types)
    active = BooleanField()


class BioSynthForm(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0)
    classes = FieldList(
        FormField(BiosynthClassForm),
        description="List of classes in this entry",
        widget=FieldListAddBtn(
            label="Add class",
        ),
    )

    modules = FieldList(
        FormField(BiosynthModuleForm),
        widget=FieldListAddBtn(label="Add module"),
    )

    operons = FieldList(
        FormField(OperonForm),
        widget=FieldListAddBtn(label="Add operon"),
    )

    paths = FieldList(
        FormField(PathForm),
        widget=FieldListAddBtn(label="Add biosynthetic pathway"),
    )

    submit = SubmitField("Save changes", widget=SubmitIndicator())
