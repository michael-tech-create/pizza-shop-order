import math 

user_lon_str = float(input().replace(",", "."))
user_lat_str = float(input().replace(",", "."))

user_lon_rad = math.radians(user_lat_str)
user_lat_rad = math.radians(user_lat_str)

closest_dist = float('inf')
closest_defib = ""

total_defib = int(input())

for i in range(total_defib):
    defib_line = input()
    fields = defib_line.split(";")

    name = fields[1]

    lon_str = fields.replace(",", ".")
    lat_str = fields.replace(",", ".")
    defib_lon = float(lon_str)
    defib_lat = float(lat_str)

    defib_lon_rad = math.radians(defib_lon)
    defib_lat_rad = math.radians(defib_lat)

    x = (defib_lon_rad - user_lon_rad) * math.cos((user_lat_rad + defib_lat_rad)/2)
    y = defib_lat_rad - user_lat_rad

    distance = math.sqrt((x * x) + (y*y)) * 6371

    if distance < closest_dist:
        closest_dist = distance

    closest_defib = name

print(closest_defib)


