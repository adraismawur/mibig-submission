from flask import Blueprint

bp_new = Blueprint("new", __file__, url_prefix="/new")

from submission.new import views
