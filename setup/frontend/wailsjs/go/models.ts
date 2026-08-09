export namespace audit {
	
	export class Check {
	    id: string;
	    label: string;
	    ok: boolean;
	    required: boolean;
	    detail: string;
	    action: string;
	
	    static createFrom(source: any = {}) {
	        return new Check(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.label = source["label"];
	        this.ok = source["ok"];
	        this.required = source["required"];
	        this.detail = source["detail"];
	        this.action = source["action"];
	    }
	}
	export class Report {
	    checks: Check[];
	    canProceed: boolean;
	    summary: string;
	    windowsOk: boolean;
	    diskOk: boolean;
	    webView2Ok: boolean;
	    ffmpegOk: boolean;
	    wingetOk: boolean;
	    nodeOk: boolean;
	    claudeOk: boolean;
	    codexOk: boolean;
	    alreadyInstalled: boolean;
	    installPath: string;
	    installedVersion: string;
	    projectsDir: string;
	
	    static createFrom(source: any = {}) {
	        return new Report(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.checks = this.convertValues(source["checks"], Check);
	        this.canProceed = source["canProceed"];
	        this.summary = source["summary"];
	        this.windowsOk = source["windowsOk"];
	        this.diskOk = source["diskOk"];
	        this.webView2Ok = source["webView2Ok"];
	        this.ffmpegOk = source["ffmpegOk"];
	        this.wingetOk = source["wingetOk"];
	        this.nodeOk = source["nodeOk"];
	        this.claudeOk = source["claudeOk"];
	        this.codexOk = source["codexOk"];
	        this.alreadyInstalled = source["alreadyInstalled"];
	        this.installPath = source["installPath"];
	        this.installedVersion = source["installedVersion"];
	        this.projectsDir = source["projectsDir"];
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

export namespace update {
	
	export class CheckResult {
	    ok: boolean;
	    installed: boolean;
	    installDir: string;
	    installedTag: string;
	    latestTag: string;
	    latestUrl: string;
	    updateAvailable: boolean;
	    message: string;
	    error?: string;
	    checkedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new CheckResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.installed = source["installed"];
	        this.installDir = source["installDir"];
	        this.installedTag = source["installedTag"];
	        this.latestTag = source["latestTag"];
	        this.latestUrl = source["latestUrl"];
	        this.updateAvailable = source["updateAvailable"];
	        this.message = source["message"];
	        this.error = source["error"];
	        this.checkedAt = source["checkedAt"];
	    }
	}

}

