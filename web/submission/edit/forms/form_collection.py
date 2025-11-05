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
    OtherForm,
    OperonMultipleForm,
)
from submission.edit.forms.biosynthesis_paths import PathMultipleForm
from submission.edit.forms.biosynthesis_modules import ModulesForm
from submission.edit.forms.tailoring import TailoringMultipleForm
from submission.edit.forms.gene_annotation import GeneAnnotationForm
from submission.edit.forms.compounds import CompoundsForm
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
    Ribosomal = RibosomalForm
    Saccharide = SaccharideForm
    Terpene = TerpeneForm
    Other = OtherForm

    operons = OperonMultipleForm
    paths = PathMultipleForm
    modules = ModulesForm

    tailoring = TailoringMultipleForm
    annotation = GeneAnnotationForm

    compounds = CompoundsForm
