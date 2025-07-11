#!/bin/bash

WORKSPACE_NAME="luggage_ctr"
WORKSPACE_DIR=$(pwd)

usage() { 
	echo "Usage: $0 [-h]" 1>&2
	echo "  Required environment variables:"
	
	if [ -z "${LUGGAGE_DB_DIAL_ADDR+x}" ]; then
		echo "    LUGGAGE_DB_DIAL_ADDR (missing)"
	else
		echo "    LUGGAGE_DB_DIAL_ADDR=$LUGGAGE_DB_DIAL_ADDR"
	fi
	if [ -z "${LUGGAGE_SERVICE_HTTP_BIND_ADDR+x}" ]; then
		echo "    LUGGAGE_SERVICE_HTTP_BIND_ADDR (missing)"
	else
		echo "    LUGGAGE_SERVICE_HTTP_BIND_ADDR=$LUGGAGE_SERVICE_HTTP_BIND_ADDR"
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


luggage_proc() {
	cd $WORKSPACE_DIR
	
	if [ -z "${LUGGAGE_DB_DIAL_ADDR+x}" ]; then
		if ! luggage_db_dial_addr; then
			return $?
		fi
	fi

	if [ -z "${LUGGAGE_SERVICE_HTTP_BIND_ADDR+x}" ]; then
		if ! luggage_service_http_bind_addr; then
			return $?
		fi
	fi

	run_luggage_proc() {
		
        cd luggage_proc
        ./luggage_proc --luggage_db.dial_addr=$LUGGAGE_DB_DIAL_ADDR --luggage_service.http.bind_addr=$LUGGAGE_SERVICE_HTTP_BIND_ADDR &
        LUGGAGE_PROC=$!
        return $?

	}

	if run_luggage_proc; then
		if [ -z "${LUGGAGE_PROC+x}" ]; then
			echo "${WORKSPACE_NAME} error starting luggage_proc: function luggage_proc did not set LUGGAGE_PROC"
			return 1
		else
			echo "${WORKSPACE_NAME} started luggage_proc"
			return 0
		fi
	else
		exitcode=$?
		echo "${WORKSPACE_NAME} aborting luggage_proc due to exitcode ${exitcode} from luggage_proc"
		return $exitcode
	fi
}


run_all() {
	echo "Running luggage_ctr"

	# Check that all necessary environment variables are set
	echo "Required environment variables:"
	missing_vars=0
	if [ -z "${LUGGAGE_DB_DIAL_ADDR+x}" ]; then
		echo "  LUGGAGE_DB_DIAL_ADDR (missing)"
		missing_vars=$((missing_vars+1))
	else
		echo "  LUGGAGE_DB_DIAL_ADDR=$LUGGAGE_DB_DIAL_ADDR"
	fi
	
	if [ -z "${LUGGAGE_SERVICE_HTTP_BIND_ADDR+x}" ]; then
		echo "  LUGGAGE_SERVICE_HTTP_BIND_ADDR (missing)"
		missing_vars=$((missing_vars+1))
	else
		echo "  LUGGAGE_SERVICE_HTTP_BIND_ADDR=$LUGGAGE_SERVICE_HTTP_BIND_ADDR"
	fi
		

	if [ "$missing_vars" -gt 0 ]; then
		echo "Aborting due to missing environment variables"
		return 1
	fi

	luggage_proc
	
	wait
}

run_all
