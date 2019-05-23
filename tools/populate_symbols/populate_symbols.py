import os, re, requests

script_dir = os.path.dirname(__file__)
add_symbol_url = "http://localhost:5000/insert_symbol_1"

def read_symbols_from_file(filename, exchange_name):
    with open(os.path.join(script_dir, filename), "r+") as f:
        lines = f.readlines()
        categories = [r.lower() for r in re.sub("\"", "", lines[0]).split(",")]
        for profile in lines[1:]:
            company = {categories[i]: re.sub("\"", "", v) for i, v in enumerate(re.sub(",\r\n", "", profile).split("\",\""))}
            if "^" in company["symbol"] or "." in company["symbol"]:
                continue
            req_data = make_request_data_from_company(company, exchange_name)
            r = requests.post(add_symbol_url, json=req_data)
            if not r.ok:
                print "company {} not added".format(company["name"])
                print r.content

def make_request_data_from_company(company, exchange_name):
   data = dict()
   data["symbol"] = company["symbol"]
   data["name"] = company["name"]
   data["market_capitalization"] = int(float(company["marketcap"]))
   data["sector"] = company["sector"] if company["sector"] != "n/a" else None
   data["industry"] = company["industry"] if company["industry"] != "n/a" else None
   data["exchange"] = exchange_name
   return data

read_symbols_from_file("symbols_data/nyse.csv", "nyse")
