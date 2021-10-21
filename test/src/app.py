from json import dumps
from uuid import uuid4
from flask import Flask, render_template, request

app = Flask(__name__)


@app.route("/")
def home():
    return "stable"


@app.route('/2.0/repositories/<w>/<r>/refs/tags')
def tags(w, r):
    return render_template("tags.json", workspace=w, repo=t)


@app.route('/2.0/repositories/<w>/<r>/refs/tags/<t>')
def tag(w, r, t):
    return render_template("tag.json", workspace=w, repo=r, tag=t)


@app.route('/2.0/repositories/<w>/<r>/pipelines/<i>')
def get_pipeline(w, r, i):
    return render_template("pipeline.json", workspace=w, repo=r, uuid=i)


@app.route('/2.0/repositories/<w>/<r>/pipelines/<i>/steps/')
def get_steps(w, r, i):
    return render_template("steps.json", workspace=w, repo=r, uuid='{' + uuid4().__str__() + '}')


@app.route('/2.0/repositories/<w>/<r>/pipelines/<pi>/steps/<si>/log')
def get_step_log(w, r, pi, si):
    return f'''
Some other logging info before outputs
--- OUTPUT JSON START ---
{{
    "string_output": "value",
    "map_output": {{
        "map_key": "value"
    }}
}}
--- OUTPUT JSON STOP ---
Some other logging info after outputs
'''


@app.route('/2.0/repositories/<w>/<r>/pipelines/', methods=['POST'])
def post_pipeline(w, r):
    data = request.get_json()
    data['UUID'] = uuid4().__str__()
    data['state'] = {
        "name": "COMPLETED",
        "type": "pipeline_state_completed",
        "result": {
            "name": "SUCCESSFUL",
            "type": "pipeline_state_completed_successful"
        }
    }
    data["build_number"] = 1
    data["created_on"] = "2021-03-10T00:32:47.890073Z"
    data["completed_on"] = "2021-03-10T00:33:08.060765Z"

    app.logger.info(dumps(data))

    return data
