export namespace config {
	
	export class BackendConfig {
	    m7_path: string;
	    better_gi_path: string;
	    zzzod_path: string;
	    genshin_auto_login_path: string;
	    ok_ww_path: string;
	    stop_on_script_error: boolean;
	    stop_on_account_error: boolean;
	    full_run_timeout_minutes: number;
	    daily_mission_timeout_minutes: number;
	    refresh_stamina_timeout_minutes: number;
	    simulated_universe_timeout_minutes: number;
	    farming_timeout_minutes: number;
	    better_gi_timeout_minutes: number;
	    better_gi_scheduler_timeout_minutes: number;
	    zzzod_timeout_minutes: number;
	    ok_ww_timeout_minutes: number;
	
	    static createFrom(source: any = {}) {
	        return new BackendConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.m7_path = source["m7_path"];
	        this.better_gi_path = source["better_gi_path"];
	        this.zzzod_path = source["zzzod_path"];
	        this.genshin_auto_login_path = source["genshin_auto_login_path"];
	        this.ok_ww_path = source["ok_ww_path"];
	        this.stop_on_script_error = source["stop_on_script_error"];
	        this.stop_on_account_error = source["stop_on_account_error"];
	        this.full_run_timeout_minutes = source["full_run_timeout_minutes"];
	        this.daily_mission_timeout_minutes = source["daily_mission_timeout_minutes"];
	        this.refresh_stamina_timeout_minutes = source["refresh_stamina_timeout_minutes"];
	        this.simulated_universe_timeout_minutes = source["simulated_universe_timeout_minutes"];
	        this.farming_timeout_minutes = source["farming_timeout_minutes"];
	        this.better_gi_timeout_minutes = source["better_gi_timeout_minutes"];
	        this.better_gi_scheduler_timeout_minutes = source["better_gi_scheduler_timeout_minutes"];
	        this.zzzod_timeout_minutes = source["zzzod_timeout_minutes"];
	        this.ok_ww_timeout_minutes = source["ok_ww_timeout_minutes"];
	    }
	}

}

