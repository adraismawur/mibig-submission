from dataclasses import dataclass
from flask import request
from submission.edit.forms.bio_synth import BioSynthForm
from submission.edit.forms.structure import StructureMultiple
from submission.edit.forms.new_entry import NewEntryForm
from submission.edit.forms.loci_tax import LociTaxonomyForm
from submission.edit.forms.biological_activity import BioActivityMultiple
from submission.edit.forms.biosynthesis import (
    NRPSForm,
    PKSForm,
    RibosomalForm,
    SaccharideForm,
    TerpeneForm,
    ClassOtherForm,
    OperonMultipleForm,
)
from submission.edit.forms.biosynthesis_paths import PathMultipleForm
from submission.edit.forms.biosynthesis_modules import (
    CalForm,
    NRPS_I_Form,
    NRPS_VI_Form,
    PKSIterativeForm,
    PKSModularForm,
    PKSTransForm,
    ModuleOtherForm,
)
from submission.edit.forms.tailoring import TailoringMultipleForm
from submission.edit.forms.gene_information import (
    AddGeneForm,
    AnnotationForm,
    DeleteGeneForm,
    GeneInformationForm,
)
from submission.edit.forms.compounds import CompoundsForm, CompoundsSubForm
from submission.edit.forms.finalize import FinalizeForm


class FormCollection:
    new = NewEntryForm
    locitax = LociTaxonomyForm
    biosynth = BioSynthForm

    structure = StructureMultiple
    bioact = BioActivityMultiple

    # Biosynthesis classes
    NRPS = NRPSForm
    PKS = PKSForm
    ribosomal = RibosomalForm
    saccharide = SaccharideForm
    terpene = TerpeneForm
    other = ClassOtherForm

    operons = OperonMultipleForm
    paths = PathMultipleForm

    # Biosynthesis modules
    cal = CalForm
    nrps_type1 = NRPS_I_Form
    nrps_type6 = NRPS_VI_Form
    pks_iterative = PKSIterativeForm
    pks_modular = PKSModularForm
    pks_trans_at = PKSTransForm

    tailoring = TailoringMultipleForm
    annotation = GeneInformationForm

    compounds = CompoundsForm
    edit_compound = CompoundsSubForm
    new_compound = CompoundsSubForm

    gene_information = GeneInformationForm

    new_addition = AddGeneForm
    new_deletion = DeleteGeneForm
    new_annotation = AnnotationForm

    edit_addition = AddGeneForm
    edit_deletion = DeleteGeneForm
    edit_annotation = AnnotationForm
