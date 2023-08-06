# Notes

```
gcloud auth application-default login
gcloud config set project <name>
```

```
gcloud compute instances list --filter="zone:us-central1-a"
```

```
c := google.client("compute", "instances")
c.list({Project: "name", Zone: "us-central1-a"})
```
