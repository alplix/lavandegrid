export namespace app {
	
	export class DiskProject {
	    URL: string;
	    DiskUsage: number;
	
	    static createFrom(source: any = {}) {
	        return new DiskProject(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.URL = source["URL"];
	        this.DiskUsage = source["DiskUsage"];
	    }
	}
	export class DiskInfo {
	    Total: number;
	    Free: number;
	    Projects: DiskProject[];
	
	    static createFrom(source: any = {}) {
	        return new DiskInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Total = source["Total"];
	        this.Free = source["Free"];
	        this.Projects = this.convertValues(source["Projects"], DiskProject);
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
	
	export class GPU {
	    Vendor: string;
	    Count: number;
	    Names: string[];
	    Driver: string;
	    Cuda: string;
	    VRAM: number;
	
	    static createFrom(source: any = {}) {
	        return new GPU(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Vendor = source["Vendor"];
	        this.Count = source["Count"];
	        this.Names = source["Names"];
	        this.Driver = source["Driver"];
	        this.Cuda = source["Cuda"];
	        this.VRAM = source["VRAM"];
	    }
	}
	export class HistPoint {
	    T: time.Time;
	    Running: number;
	
	    static createFrom(source: any = {}) {
	        return new HistPoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.T = this.convertValues(source["T"], time.Time);
	        this.Running = source["Running"];
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
	export class HostSpec {
	    OS: string;
	    OSVersion: string;
	    CPU: string;
	    Cores: number;
	    Flops: number;
	    Memory: number;
	    DiskFree: number;
	    DiskTotal: number;
	    CPID: string;
	    GPUs: GPU[];
	
	    static createFrom(source: any = {}) {
	        return new HostSpec(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.OS = source["OS"];
	        this.OSVersion = source["OSVersion"];
	        this.CPU = source["CPU"];
	        this.Cores = source["Cores"];
	        this.Flops = source["Flops"];
	        this.Memory = source["Memory"];
	        this.DiskFree = source["DiskFree"];
	        this.DiskTotal = source["DiskTotal"];
	        this.CPID = source["CPID"];
	        this.GPUs = this.convertValues(source["GPUs"], GPU);
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
export class MsgLine {
	    Seq: number;
	    Pri: number;
	    Time: number;
	    Body: string;
	    Project: string;

	    static createFrom(source: any = {}) {
	        return new MsgLine(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Seq = source["Seq"];
	        this.Pri = source["Pri"];
	        this.Time = source["Time"];
	        this.Body = source["Body"];
	        this.Project = source["Project"];
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
	export class ProjectInfo {
	    Name: string;
	    URL: string;
	    Venue: string;
	    UserName: string;
	    TeamName: string;
	    UserCredit: number;
	    RAC: number;
	    HostCredit: number;
	    HostRAC: number;
	    Share: number;
	    Suspended: boolean;
	    NoMoreWork: boolean;
	    Pending: boolean;
	    Ended: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ProjectInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.URL = source["URL"];
	        this.Venue = source["Venue"];
	        this.UserName = source["UserName"];
	        this.TeamName = source["TeamName"];
	        this.UserCredit = source["UserCredit"];
	        this.RAC = source["RAC"];
	        this.HostCredit = source["HostCredit"];
	        this.HostRAC = source["HostRAC"];
	        this.Share = source["Share"];
	        this.Suspended = source["Suspended"];
	        this.NoMoreWork = source["NoMoreWork"];
	        this.Pending = source["Pending"];
	        this.Ended = source["Ended"];
	    }
	}
	export class Totals {
	    Running: number;
	    Paused: number;
	    Queued: number;
	    Errors: number;
	    Downloads: number;
	    Uploads: number;
	    Memory: number;
	    Credit: number;
	    RAC: number;
	
	    static createFrom(source: any = {}) {
	        return new Totals(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Running = source["Running"];
	        this.Paused = source["Paused"];
	        this.Queued = source["Queued"];
	        this.Errors = source["Errors"];
	        this.Downloads = source["Downloads"];
	        this.Uploads = source["Uploads"];
	        this.Memory = source["Memory"];
	        this.Credit = source["Credit"];
	        this.RAC = source["RAC"];
	    }
	}
	export class Transfer {
	    Name: string;
	    URL: string;
	    ProjectName: string;
	    Upload: boolean;
	    Total: number;
	    Done: number;
	    Progress: number;
	    Paused: boolean;
	    Finished: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Transfer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.URL = source["URL"];
	        this.ProjectName = source["ProjectName"];
	        this.Upload = source["Upload"];
	        this.Total = source["Total"];
	        this.Done = source["Done"];
	        this.Progress = source["Progress"];
	        this.Paused = source["Paused"];
	        this.Finished = source["Finished"];
	    }
	}
	export class Task {
	    Name: string;
	    Wu: string;
	    URL: string;
	    ProjectName: string;
	    Status: string;
	    Progress: number;
	    Elapsed: number;
	    CPUTime: number;
	    ETA: number;
	    Deadline: number;
	    CPTime: number;
	    Exit: number;
	    Mem: number;
	    Resources: string;
	    Slot: number;
	    AppVersion: string;
	    Active: boolean;
	    Suspended: boolean;
	    Ready: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Wu = source["Wu"];
	        this.URL = source["URL"];
	        this.ProjectName = source["ProjectName"];
	        this.Status = source["Status"];
	        this.Progress = source["Progress"];
	        this.Elapsed = source["Elapsed"];
	        this.CPUTime = source["CPUTime"];
	        this.ETA = source["ETA"];
	        this.Deadline = source["Deadline"];
	        this.CPTime = source["CPTime"];
	        this.Exit = source["Exit"];
	        this.Mem = source["Mem"];
	        this.Resources = source["Resources"];
	        this.Slot = source["Slot"];
	        this.AppVersion = source["AppVersion"];
	        this.Active = source["Active"];
	        this.Suspended = source["Suspended"];
	        this.Ready = source["Ready"];
	    }
	}
	export class Snapshot {
	    HostID: string;
	    Demo: boolean;
	    Online: boolean;
	    Error: string;
	    Version: string;
	    TS: time.Time;
	    HostInfo: HostSpec;
	    Projects: ProjectInfo[];
	    Tasks: Task[];
	    Transfers: Transfer[];
	    Messages: MsgLine[];
	    Totals: Totals;
	    TaskMode: string;
	    NetMode: string;
	
	    static createFrom(source: any = {}) {
	        return new Snapshot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.HostID = source["HostID"];
	        this.Demo = source["Demo"];
	        this.Online = source["Online"];
	        this.Error = source["Error"];
	        this.Version = source["Version"];
	        this.TS = this.convertValues(source["TS"], time.Time);
	        this.HostInfo = this.convertValues(source["HostInfo"], HostSpec);
	        this.Projects = this.convertValues(source["Projects"], ProjectInfo);
	        this.Tasks = this.convertValues(source["Tasks"], Task);
	        this.Transfers = this.convertValues(source["Transfers"], Transfer);
	        this.Messages = this.convertValues(source["Messages"], MsgLine);
	        this.Totals = this.convertValues(source["Totals"], Totals);
	        this.TaskMode = source["TaskMode"];
	        this.NetMode = source["NetMode"];
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
	export class StatPoint {
	    Day: string;
	    HostCredit: number;
	    UserCredit: number;
	
	    static createFrom(source: any = {}) {
	        return new StatPoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Day = source["Day"];
	        this.HostCredit = source["HostCredit"];
	        this.UserCredit = source["UserCredit"];
	    }
	}
	export class StatSeries {
	    URL: string;
	    Name: string;
	    Daily: StatPoint[];
	
	    static createFrom(source: any = {}) {
	        return new StatSeries(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.URL = source["URL"];
	        this.Name = source["Name"];
	        this.Daily = this.convertValues(source["Daily"], StatPoint);
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
	
	
	
	export class XferPoint {
	    When: number;
	    Up: number;
	    Down: number;
	
	    static createFrom(source: any = {}) {
	        return new XferPoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.When = source["When"];
	        this.Up = source["Up"];
	        this.Down = source["Down"];
	    }
	}

}

export namespace main {
	
	export class HostCfgJSON {
	    id: string;
	    name: string;
	    host: string;
	    port: number;
	    password: string;
	    demo: boolean;
	
	    static createFrom(source: any = {}) {
	        return new HostCfgJSON(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.password = source["password"];
	        this.demo = source["demo"];
	    }
	}

}

export namespace time {
	
	export class Time {
	
	
	    static createFrom(source: any = {}) {
	        return new Time(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

