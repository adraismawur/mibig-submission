from flask import Blueprint

bp_diff = Blueprint("diff", __file__, url_prefix="/diff")

from submission.diff import views
