#!/bin/bash

WORKSPACE_NAME="reserv_ctr"
WORKSPACE_DIR=$(pwd)

usage() { 
	echo "Usage: $0 [-h]" 1>&2
	echo "  Required environment variables:"
	
	if [ -z "${RESERV_DB_DIAL_ADDR+x}" ]; then
		echo "    RESERV_DB_DIAL_ADDR (missing)"
	else
		echo "    RESERV_DB_DIAL_ADDR=$RESERV_DB_DIAL_ADDR"
	fi
	if [ -z "${RESERV_SERVICE_HTTP_BIND_ADDR+x}" ]; then
		echo "    RESERV_SERVICE_HTTP_BIND_ADDR (missing)"
	else
		echo "    RESERV_SERVICE_HTTP_BIND_ADDR=$RESERV_SERVICE_HTTP_BIND_ADDR"
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


reserv_proc() {
	cd $WORKSPACE_DIR
	
	if [ -z "${RESERV_DB_DIAL_ADDR+x}" ]; then
		if ! reserv_db_dial_addr; then
			return $?
		fi
	fi

	if [ -z "${RESERV_SERVICE_HTTP_BIND_ADDR+x}" ]; then
		if ! reserv_service_http_bind_addr; then
			return $?
		fi
	fi

	run_reserv_proc() {
		
        cd reserv_proc
        ./reserv_proc --reserv_db.dial_addr=$RESERV_DB_DIAL_ADDR --reserv_service.http.bind_addr=$RESERV_SERVICE_HTTP_BIND_ADDR &
        RESERV_PROC=$!
        return $?

	}

	if run_reserv_proc; then
		if [ -z "${RESERV_PROC+x}" ]; then
			echo "${WORKSPACE_NAME} error starting reserv_proc: function reserv_proc did not set RESERV_PROC"
			return 1
		else
			echo "${WORKSPACE_NAME} started reserv_proc"
			return 0
		fi
	else
		exitcode=$?
		echo "${WORKSPACE_NAME} aborting reserv_proc due to exitcode ${exitcode} from reserv_proc"
		return $exitcode
	fi
}


run_all() {
	echo "Running reserv_ctr"

	# Check that all necessary environment variables are set
	echo "Required environment variables:"
	missing_vars=0
	if [ -z "${RESERV_DB_DIAL_ADDR+x}" ]; then
		echo "  RESERV_DB_DIAL_ADDR (missing)"
		missing_vars=$((missing_vars+1))
	else
		echo "  RESERV_DB_DIAL_ADDR=$RESERV_DB_DIAL_ADDR"
	fi
	
	if [ -z "${RESERV_SERVICE_HTTP_BIND_ADDR+x}" ]; then
		echo "  RESERV_SERVICE_HTTP_BIND_ADDR (missing)"
		missing_vars=$((missing_vars+1))
	else
		echo "  RESERV_SERVICE_HTTP_BIND_ADDR=$RESERV_SERVICE_HTTP_BIND_ADDR"
	fi
		

	if [ "$missing_vars" -gt 0 ]; then
		echo "Aborting due to missing environment variables"
		return 1
	fi

	reserv_proc
	
	wait
}

run_all
