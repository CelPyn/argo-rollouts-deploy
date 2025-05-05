FOLDERS="01-recreate 02-rolling 03-vanilla-blue-green 04-basic-rollout 05-alb-rollout"

for dest in $FOLDERS; do
  echo "Applying $dest"
  kubectl apply -f "$dest"
  echo "-----------------------------------------------------------------------"
done
