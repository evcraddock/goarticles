#!/usr/bin/env bash
helm repo add stable https://kubernetes-charts.storage.googleapis.com/
helm install --name goarticles-db -f values.yaml stable/mongodb-replicaset