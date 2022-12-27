#!/bin/bash


# Variables
dashboard=./dashboard_energy.json
dashboard_final=./dashboard_energy_final.json
config=./energy-rates.conf
file_suffix=`date +%s | sha256sum | base64 | head -c 6`

# Check Target File
if [ ! -f "$config" ]; then
    echo "Config does not exist at $config"
    exit 0;
fi



# Dashboard Variable Replace Logic
if [ -f "$dashboard" ]; then

    # Prepare Config
    cp $config $config-$file_suffix
    echo -e "\n" >> $config-$file_suffix
    # Prepare Dashboard
    cp $dashboard $dashboard-$file_suffix

    while IFS='=' read -r key value
        do
            [ -z "$value" ] && continue
            
            replace_str="s/@$key@/$value/"
            echo "Using $value for $key Replace String $replace_str"
            sed -i $replace_str $dashboard-$file_suffix
    done < $config-$file_suffix


    # Remove Temp Files
    mv $dashboard-$file_suffix $dashboard_final
    rm $config-$file_suffix

else 
    echo "$dashboard does not exist."
fi

# Run Containers
docker-compose down && docker-compose up

