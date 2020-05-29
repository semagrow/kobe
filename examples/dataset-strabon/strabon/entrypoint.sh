# Exit at any error
set -e

# Update the given property file by overriding its values based on the environment's ones
#
# To override an existing property file's property, environment key has to be named as follow: <PREFIX>_<NAME>, where:
#   - <PREFIX> is the given environment key's prefix
#   - <NAME> is the property file's property key name
#
# @param $1 path to the property file
# @param $2 the environment key's prefix
function updatePropertyFile {
    local propertyFilePath=$1
    local environmentPrefix=$2

    local updatedPropertyFile=''
    # For each property file's property, then
    while read -r property || [ -n "$property" ]; do
        if [ -z "$property" ]; then
            continue
        fi;
        # 1. Extract property's key, property's value and potentially associated environment's value
        local propertyKey=`echo "$property" | grep -o -P '^[^=]+'`
        local propertyValue=`echo "$property" | grep -o -P '=.*$' | sed 's/^=//'`
        local overriddenPropertyValue=`printenv "${environmentPrefix}_${propertyKey}"`
        # 2. Iteratively construct the new property file
        if [ -n "$updatedPropertyFile" ]; then
                updatedPropertyFile="$updatedPropertyFile"'\n'
        fi
        updatedPropertyFile="${updatedPropertyFile}${propertyKey}=${overriddenPropertyValue:-$propertyValue}"
    done < $propertyFilePath

    # Finally override the property file by the new constructed one
    echo -e "$updatedPropertyFile" > $propertyFilePath
}

# Enhance Strabon's property files
function enhancePropertyFiles {
    updatePropertyFile $STRABON_HOME/WEB-INF/credentials.properties STRABON_CREDENTIALS
    updatePropertyFile $STRABON_HOME/WEB-INF/connection.properties STRABON_CONNECTION
}

# Run next command (Docker's CMD)
function runNext {
    exec "$@"
}

# Main entry point
function main {
    enhancePropertyFiles
    runNext "$@"
}

main "$@"