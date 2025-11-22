export namespace main {
	
	export class App {
	
	
	    static Greet(arg1:string):Promise<string>;
	
	    static SaveMessage(arg1:string,arg2:string,arg3:string):Promise<void>;
	
	    static GetMessageHistory(arg1:string,arg2:number):Promise<any>;
	
	    static SavePeer(arg1:string,arg2:string,arg3:string):Promise<void>;
	
	    static GetPeers():Promise<any>;
	
	    static ShowNotification(arg1:string,arg2:string):Promise<void>;
	
	}

}
