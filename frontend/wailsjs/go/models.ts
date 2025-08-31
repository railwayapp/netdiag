export namespace core {
	
	export enum DiagnosticUpdateType {
	    START = "start",
	    STEP_START = "step_start",
	    STEP_PROGRESS = "step_progress",
	    DONE = "done",
	    ERROR = "error",
	}
	export class DiagnosticUpdate {
	    type: DiagnosticUpdateType;
	    message: string;
	    data: string;
	
	    static createFrom(source: any = {}) {
	        return new DiagnosticUpdate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.message = source["message"];
	        this.data = source["data"];
	    }
	}

}

