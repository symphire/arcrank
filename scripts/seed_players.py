import argparse
import json
import os
import random
import requests
import sys

BASE_URL = "http://localhost:8080"
USER_COUNT = 10_000

regions = ["CN", "EU", "JP", "NA", "KR"]
classes = ["Mage", "Warrior", "Rogue", "Cleric"]
adjectives = ["dark", "swift", "silent", "brave", "iron"]
nouns = ["blade", "wolf", "mage", "knight", "arrow"]

def random_username(i):
    return f"{random.choice(adjectives)}_{random.choice(nouns)}_{i}"

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--out-dir",
        default="./k6/scripts",
        help="Directory to write ids.json, relative to current working directory (default: ./k6/scripts)",
    )
    parser.add_argument(
        "--sleep",
        type=float,
        default=0.01,
        help="Sleep time (seconds) between requests (default: 0.01)",
    )
    args = parser.parse_args()

    out_dir = os.path.abspath(os.path.expanduser(args.out_dir))
    os.makedirs(out_dir, exist_ok=True)
    out_path = os.path.join(out_dir, "ids.json")

    ids = []
    err = False

    for i in range(USER_COUNT):
        payload = {
            "username": random_username(i),
            "region": random.choice(regions),
            "class": random.choice(classes),
        }

        resp = requests.post(
            f"{BASE_URL}/players/",
            json=payload,
            timeout=5,
        )

        if resp.status_code != 201:
            print("seed failed:", resp.status_code, resp.text, file=sys.stderr)
            err = True
            break

        body = resp.json()
        if "id" not in body:
            print("no id returned:", body, file=sys.stderr)
            err = True
            break

        ids.append(body["id"])
        # print("created:", body["id"])

    with open(out_path, "w") as f:
        json.dump(ids, f, indent=2)

    print(f"seeded {len(ids)} users")
    if err:
        sys.exit(1)

if __name__ == "__main__":
    main()
