#!/bin/bash

WORKSPACE_NAME="search_ctr"
WORKSPACE_DIR=$(pwd)

usage() { 
	echo "Usage: $0 [-h]" 1>&2
	echo "  Required environment variables:"
	
	if [ -z "${LUGGAGE_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		echo "    LUGGAGE_SERVICE_HTTP_DIAL_ADDR (missing)"
	else
		echo "    LUGGAGE_SERVICE_HTTP_DIAL_ADDR=$LUGGAGE_SERVICE_HTTP_DIAL_ADDR"
	fi
	if [ -z "${RESERV_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		echo "    RESERV_SERVICE_HTTP_DIAL_ADDR (missing)"
	else
		echo "    RESERV_SERVICE_HTTP_DIAL_ADDR=$RESERV_SERVICE_HTTP_DIAL_ADDR"
	fi
	if [ -z "${SEARCH_SERVICE_HTTP_BIND_ADDR+x}" ]; then
		echo "    SEARCH_SERVICE_HTTP_BIND_ADDR (missing)"
	else
		echo "    SEARCH_SERVICE_HTTP_BIND_ADDR=$SEARCH_SERVICE_HTTP_BIND_ADDR"
	fi
		
	exit 1; 
}

while getopts "h" flag; do
	case $flag in
		*)
		usage
		;;
	esac
done


search_proc() {
	cd $WORKSPACE_DIR
	
	if [ -z "${RESERV_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		if ! reserv_service_http_dial_addr; then
			return $?
		fi
	fi

	if [ -z "${LUGGAGE_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		if ! luggage_service_http_dial_addr; then
			return $?
		fi
	fi

	if [ -z "${SEARCH_SERVICE_HTTP_BIND_ADDR+x}" ]; then
		if ! search_service_http_bind_addr; then
			return $?
		fi
	fi

	run_search_proc() {
		
        cd search_proc
        ./search_proc --reserv_service.http.dial_addr=$RESERV_SERVICE_HTTP_DIAL_ADDR --luggage_service.http.dial_addr=$LUGGAGE_SERVICE_HTTP_DIAL_ADDR --search_service.http.bind_addr=$SEARCH_SERVICE_HTTP_BIND_ADDR &
        SEARCH_PROC=$!
        return $?

	}

	if run_search_proc; then
		if [ -z "${SEARCH_PROC+x}" ]; then
			echo "${WORKSPACE_NAME} error starting search_proc: function search_proc did not set SEARCH_PROC"
			return 1
		else
			echo "${WORKSPACE_NAME} started search_proc"
			return 0
		fi
	else
		exitcode=$?
		echo "${WORKSPACE_NAME} aborting search_proc due to exitcode ${exitcode} from search_proc"
		return $exitcode
	fi
}


run_all() {
	echo "Running search_ctr"

	# Check that all necessary environment variables are set
	echo "Required environment variables:"
	missing_vars=0
	if [ -z "${LUGGAGE_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		echo "  LUGGAGE_SERVICE_HTTP_DIAL_ADDR (missing)"
		missing_vars=$((missing_vars+1))
	else
		echo "  LUGGAGE_SERVICE_HTTP_DIAL_ADDR=$LUGGAGE_SERVICE_HTTP_DIAL_ADDR"
	fi
	
	if [ -z "${RESERV_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		echo "  RESERV_SERVICE_HTTP_DIAL_ADDR (missing)"
		missing_vars=$((missing_vars+1))
	else
		echo "  RESERV_SERVICE_HTTP_DIAL_ADDR=$RESERV_SERVICE_HTTP_DIAL_ADDR"
	fi
	
	if [ -z "${SEARCH_SERVICE_HTTP_BIND_ADDR+x}" ]; then
		echo "  SEARCH_SERVICE_HTTP_BIND_ADDR (missing)"
		missing_vars=$((missing_vars+1))
	else
		echo "  SEARCH_SERVICE_HTTP_BIND_ADDR=$SEARCH_SERVICE_HTTP_BIND_ADDR"
	fi
		

	if [ "$missing_vars" -gt 0 ]; then
		echo "Aborting due to missing environment variables"
		return 1
	fi

	search_proc
	
	wait
}

run_all
