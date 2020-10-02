import requests
import time
import json

def get_session_id() -> str:
    r = requests.get("https://tableau.gatech.edu/t/GT/views/GTCovid-19Tracking/GeorgiaInstituteofTechnologyCovid-19Data?%3Aembed=y")
    content = r.text
    sid_index = content.find("sessionid")
    sid_hash_index = content.find("sessionIdHash")
    c_hash_string = content[sid_index+27:sid_hash_index]
    c_hash_string = c_hash_string.replace("&quot;", "")
    c_hash_string = c_hash_string.replace("&#x3a;", ":")[:-1]
    return c_hash_string


short_name_conversions = {
    "Past Seven Days Rolling Averages": "sevenDayRA",
    "Testing Count since August 2020": "testingSinceAug",
    "Student Test Count": "studentTestCount",
    "Employee/Affiliate Test Count": "facultyTestCount",
    "Isolation/Quarantine In Use Bed Count": "isolationInUse",
    "Total Beds": "totalBeds",
    "OverTime Chart": "otChart",
    "Number of Tests Per Day": "numTestsDaily",
    "Chart Title": "chartTitle",
    "Count of Positive Cases since March 2020": "positiveSinceMar",
    "Student Count of Positive Cases": "studentPositiveCount",
    "Employee/Affiliate Count of Positive Cases": "facultyPositiveCount",
    "OverTime Table": "otTable"
}


# get the sheet with the graph and other info
def getFirstSheet(session_id: str) -> dict:
    url = f"https://tableau.gatech.edu/vizql/t/GT/w/GTCovid-19Tracking/v/GeorgiaInstituteofTechnologyCovid-19Data/bootstrapSession/sessions/{session_id}"
    payload = "sheet_id=Georgia%252520Institute%252520of%252520Technology%252520Covid-19%252520Data"
    headers = {
        'user-agent': "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4181.8 Safari/537.36",
        'content-type': "application/x-www-form-urlencoded",
        'cookie': "tableau_locale=en"
    }

    r = requests.post(url, headers=headers, data=payload)
    content = r.text
    content = content[content.find(";")+1:content.find("20;{\"")]

    resp_dict = json.loads(content)
    zones = resp_dict["worldUpdate"]["applicationPresModel"]["workbookPresModel"]["dashboardPresModel"]["zones"]

    ret_dict = {}
    for k,v in zones.items():
        try:
            if "png" in v["presModelHolder"]["visual"]["cacheUrlInfoJson"]:
                path = json.loads(v["presModelHolder"]["visual"]["cacheUrlInfoJson"].replace("%SESSIONID%", session_id))["url"]
                url = "https://tableau.gatech.edu" + path
                ret_dict.update({short_name_conversions[v["worksheet"]]: url})
        except KeyError:
            continue
    return ret_dict

def getSecondSheet(session_id: str) -> dict:
    url = f"https://tableau.gatech.edu/vizql/t/GT/w/GTCovid-19Tracking/v/GeorgiaInstituteofTechnologyCovid-19Data/sessions/{session_id}/commands/tabsrv/ensure-layout-for-sheet"

    payload = "-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"targetSheet\"\r\n\r\nGeorgia Institute of Technology Covid-19 Data \r\n-----011000010111000001101001--\r\n"
    headers = {
            'cookie': "tableau_locale=en",
            'user-agent': "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4181.8 Safari/537.36",
            'content-type': "multipart/form-data; boundary=---011000010111000001101001"
            }

    response = requests.request("POST", url, data=payload, headers=headers)

    zones = response.json()["vqlCmdResponse"]["layoutStatus"]["applicationPresModel"]["workbookPresModel"]["dashboardPresModel"]["zones"]

    ret_dict = {}
    for k,v in zones.items():
        try:
            if "png" in v["presModelHolder"]["visual"]["cacheUrlInfoJson"]:
                path = json.loads(v["presModelHolder"]["visual"]["cacheUrlInfoJson"].replace("%SESSIONID%", session_id))["url"]
                url = "https://tableau.gatech.edu" + path
                ret_dict.update({short_name_conversions[v["worksheet"]]: url})
        except KeyError:
            continue
    return ret_dict

sid = get_session_id()
print(sid)
print(getFirstSheet(sid))
print("\n\n\n\n")
print(getSecondSheet(sid))


