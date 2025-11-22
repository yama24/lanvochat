export namespace database {
	
	export class Peer {
	    id: number;
	    peer_id: string;
	    name: string;
	    ip_address: string;
	    // Go type: time
	    last_seen: any;
	    is_online: boolean;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Peer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.peer_id = source["peer_id"];
	        this.name = source["name"];
	        this.ip_address = source["ip_address"];
	        this.last_seen = this.convertValues(source["last_seen"], null);
	        this.is_online = source["is_online"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

