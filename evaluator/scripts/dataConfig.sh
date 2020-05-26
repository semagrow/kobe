# !/bin/sh

echo '@prefix fluid: <http://fluidops.org/config#> .'
echo ''
echo '<http://endpoint> fluid:store "SPARQLEndpoint" ;'
echo '                  fluid:SPARQLEndpoint "'$1'" .'

