import pandas

json_path = input("Path to jsonfile: ")
xlsx_path = input("New xlsx filename: ")

if not json_path.endswith(".json"):
	exit("Первый файл должен быть с расширением .json !")
if not xlsx_path.endswith(".xlsx"):
	exit("Второй файл должен быть с расширением .xlsx !")

pandas.read_json(json_path).to_excel(xlsx_path, index=False, engine='xlsxwriter')
print("Success.")