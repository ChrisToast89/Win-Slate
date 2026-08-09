export namespace types {
	
	export class ArtDeptSheet {
	    id: string;
	    kind: string;
	    name: string;
	    description: string;
	    materials: string;
	    condition: string;
	    era: string;
	    distinctive: string;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new ArtDeptSheet(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.kind = source["kind"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.materials = source["materials"];
	        this.condition = source["condition"];
	        this.era = source["era"];
	        this.distinctive = source["distinctive"];
	        this.notes = source["notes"];
	    }
	}
	export class AudioFingerprint {
	    durationSec: number;
	    sampledSec: number;
	    bpmEstimate?: number;
	    bpmConfidence: string;
	    pitchMedianHz?: number;
	    pitchSpreadSemitones?: number;
	    voicedRatio: number;
	    brightness: string;
	    dynamicRangeDb: number;
	    energyArc: string;
	    silenceRatio: number;
	    longestSilenceSec: number;
	
	    static createFrom(source: any = {}) {
	        return new AudioFingerprint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.durationSec = source["durationSec"];
	        this.sampledSec = source["sampledSec"];
	        this.bpmEstimate = source["bpmEstimate"];
	        this.bpmConfidence = source["bpmConfidence"];
	        this.pitchMedianHz = source["pitchMedianHz"];
	        this.pitchSpreadSemitones = source["pitchSpreadSemitones"];
	        this.voicedRatio = source["voicedRatio"];
	        this.brightness = source["brightness"];
	        this.dynamicRangeDb = source["dynamicRangeDb"];
	        this.energyArc = source["energyArc"];
	        this.silenceRatio = source["silenceRatio"];
	        this.longestSilenceSec = source["longestSilenceSec"];
	    }
	}
	export class BackendStatus {
	    available: boolean;
	    version?: string;
	
	    static createFrom(source: any = {}) {
	        return new BackendStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available = source["available"];
	        this.version = source["version"];
	    }
	}
	export class BeatDirection {
	    from: number;
	    to: number;
	    text: string;
	
	    static createFrom(source: any = {}) {
	        return new BeatDirection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.from = source["from"];
	        this.to = source["to"];
	        this.text = source["text"];
	    }
	}
	export class BrainRequest {
	    id: string;
	    task: string;
	    system: string;
	    prompt: string;
	    images?: string[];
	    tier: string;
	    expectJson?: boolean;
	    localEndpoint?: string;
	    localModel?: string;
	    backend?: string;
	
	    static createFrom(source: any = {}) {
	        return new BrainRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.task = source["task"];
	        this.system = source["system"];
	        this.prompt = source["prompt"];
	        this.images = source["images"];
	        this.tier = source["tier"];
	        this.expectJson = source["expectJson"];
	        this.localEndpoint = source["localEndpoint"];
	        this.localModel = source["localModel"];
	        this.backend = source["backend"];
	    }
	}
	export class BrainResult {
	    id: string;
	    ok: boolean;
	    text: string;
	    json?: any;
	    error?: string;
	    elapsedMs: number;
	
	    static createFrom(source: any = {}) {
	        return new BrainResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.ok = source["ok"];
	        this.text = source["text"];
	        this.json = source["json"];
	        this.error = source["error"];
	        this.elapsedMs = source["elapsedMs"];
	    }
	}
	export class LocalBackendStatus {
	    available: boolean;
	    version?: string;
	    endpoint?: string;
	
	    static createFrom(source: any = {}) {
	        return new LocalBackendStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available = source["available"];
	        this.version = source["version"];
	        this.endpoint = source["endpoint"];
	    }
	}
	export class BrainStatus {
	    claude: BackendStatus;
	    codex: BackendStatus;
	    local: LocalBackendStatus;
	
	    static createFrom(source: any = {}) {
	        return new BrainStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.claude = this.convertValues(source["claude"], BackendStatus);
	        this.codex = this.convertValues(source["codex"], BackendStatus);
	        this.local = this.convertValues(source["local"], LocalBackendStatus);
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
	export class CharacterSheet {
	    id: string;
	    name: string;
	    age: string;
	    gender: string;
	    ethnicity: string;
	    faceFeatures: string;
	    hair: string;
	    clothing: string;
	    expression: string;
	    eyeDirection: string;
	    mood: string;
	    environment: string;
	    keyLightSide: string;
	    lightingMood: string;
	    scenario: string;
	    notes: string;
	    images?: string[];
	
	    static createFrom(source: any = {}) {
	        return new CharacterSheet(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.age = source["age"];
	        this.gender = source["gender"];
	        this.ethnicity = source["ethnicity"];
	        this.faceFeatures = source["faceFeatures"];
	        this.hair = source["hair"];
	        this.clothing = source["clothing"];
	        this.expression = source["expression"];
	        this.eyeDirection = source["eyeDirection"];
	        this.mood = source["mood"];
	        this.environment = source["environment"];
	        this.keyLightSide = source["keyLightSide"];
	        this.lightingMood = source["lightingMood"];
	        this.scenario = source["scenario"];
	        this.notes = source["notes"];
	        this.images = source["images"];
	    }
	}
	export class ChatMsg {
	    role: string;
	    text: string;
	    receipts?: string[];
	
	    static createFrom(source: any = {}) {
	        return new ChatMsg(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.role = source["role"];
	        this.text = source["text"];
	        this.receipts = source["receipts"];
	    }
	}
	export class CircledTake {
	    project: string;
	    shot?: string;
	    mediaPath: string;
	    fileName: string;
	    rating: number;
	    inSec?: number;
	    outSec?: number;
	
	    static createFrom(source: any = {}) {
	        return new CircledTake(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.project = source["project"];
	        this.shot = source["shot"];
	        this.mediaPath = source["mediaPath"];
	        this.fileName = source["fileName"];
	        this.rating = source["rating"];
	        this.inSec = source["inSec"];
	        this.outSec = source["outSec"];
	    }
	}
	export class CustomSetup {
	    id: string;
	    label: string;
	    snippet: string;
	    section: string;
	    tags: string[];
	    favorite: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CustomSetup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.label = source["label"];
	        this.snippet = source["snippet"];
	        this.section = source["section"];
	        this.tags = source["tags"];
	        this.favorite = source["favorite"];
	    }
	}
	export class ElementSheet {
	    lensing: string;
	    lighting: string;
	    palette: string;
	    composition: string;
	    movement: string;
	    texture: string;
	    mood: string;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new ElementSheet(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lensing = source["lensing"];
	        this.lighting = source["lighting"];
	        this.palette = source["palette"];
	        this.composition = source["composition"];
	        this.movement = source["movement"];
	        this.texture = source["texture"];
	        this.mood = source["mood"];
	        this.notes = source["notes"];
	    }
	}
	
	export class LocalModelInfo {
	    id: string;
	
	    static createFrom(source: any = {}) {
	        return new LocalModelInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	    }
	}
	export class LocalModelsResult {
	    endpoint?: string;
	    models: LocalModelInfo[];
	
	    static createFrom(source: any = {}) {
	        return new LocalModelsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.endpoint = source["endpoint"];
	        this.models = this.convertValues(source["models"], LocalModelInfo);
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
	export class LocalOpts {
	    endpoint?: string;
	    model?: string;
	
	    static createFrom(source: any = {}) {
	        return new LocalOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.endpoint = source["endpoint"];
	        this.model = source["model"];
	    }
	}
	export class LocationSheet {
	    id: string;
	    name: string;
	    interiorExterior: string;
	    description: string;
	    timeOfDay: string;
	    weather: string;
	    architecture: string;
	    textures: string;
	    practicalLights: string;
	    notes: string;
	    images?: string[];
	
	    static createFrom(source: any = {}) {
	        return new LocationSheet(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.interiorExterior = source["interiorExterior"];
	        this.description = source["description"];
	        this.timeOfDay = source["timeOfDay"];
	        this.weather = source["weather"];
	        this.architecture = source["architecture"];
	        this.textures = source["textures"];
	        this.practicalLights = source["practicalLights"];
	        this.notes = source["notes"];
	        this.images = source["images"];
	    }
	}
	export class MediaIngestResult {
	    kind: string;
	    frames: string[];
	
	    static createFrom(source: any = {}) {
	        return new MediaIngestResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.kind = source["kind"];
	        this.frames = source["frames"];
	    }
	}
	export class MusicCue {
	    id: string;
	    name: string;
	    sceneRef: string;
	    intent: string;
	    genre: string;
	    mood: string;
	    tempo: string;
	    instrumentation: string;
	    era: string;
	    structure: string;
	    vocals: string;
	    lyricTheme: string;
	    lyrics: string;
	    durationSec?: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new MusicCue(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.sceneRef = source["sceneRef"];
	        this.intent = source["intent"];
	        this.genre = source["genre"];
	        this.mood = source["mood"];
	        this.tempo = source["tempo"];
	        this.instrumentation = source["instrumentation"];
	        this.era = source["era"];
	        this.structure = source["structure"];
	        this.vocals = source["vocals"];
	        this.lyricTheme = source["lyricTheme"];
	        this.lyrics = source["lyrics"];
	        this.durationSec = source["durationSec"];
	        this.notes = source["notes"];
	    }
	}
	export class VoiceSheet {
	    id: string;
	    name: string;
	    characterId?: string;
	    ageGender: string;
	    accent: string;
	    timbre: string;
	    pitch: string;
	    pacing: string;
	    energy: string;
	    texture: string;
	    emotionalRange: string;
	    sampleLine: string;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new VoiceSheet(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.characterId = source["characterId"];
	        this.ageGender = source["ageGender"];
	        this.accent = source["accent"];
	        this.timbre = source["timbre"];
	        this.pitch = source["pitch"];
	        this.pacing = source["pacing"];
	        this.energy = source["energy"];
	        this.texture = source["texture"];
	        this.emotionalRange = source["emotionalRange"];
	        this.sampleLine = source["sampleLine"];
	        this.notes = source["notes"];
	    }
	}
	export class Reference {
	    id: string;
	    path: string;
	    kind: string;
	    label: string;
	    frames: string[];
	    elements?: ElementSheet;
	    addedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Reference(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.path = source["path"];
	        this.kind = source["kind"];
	        this.label = source["label"];
	        this.frames = source["frames"];
	        this.elements = this.convertValues(source["elements"], ElementSheet);
	        this.addedAt = source["addedAt"];
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
	export class StyleProfile {
	    id: string;
	    source: string;
	    kind: string;
	    tone: string;
	    palette: string;
	    lighting: string;
	    lensLanguage: string;
	    movement: string;
	    blocking: string;
	    editorial: string;
	    notes: string;
	    images?: string[];
	
	    static createFrom(source: any = {}) {
	        return new StyleProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.source = source["source"];
	        this.kind = source["kind"];
	        this.tone = source["tone"];
	        this.palette = source["palette"];
	        this.lighting = source["lighting"];
	        this.lensLanguage = source["lensLanguage"];
	        this.movement = source["movement"];
	        this.blocking = source["blocking"];
	        this.editorial = source["editorial"];
	        this.notes = source["notes"];
	        this.images = source["images"];
	    }
	}
	export class Take {
	    id: string;
	    loggedAt: string;
	    model: string;
	    prompt: string;
	    rating: string;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new Take(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.loggedAt = source["loggedAt"];
	        this.model = source["model"];
	        this.prompt = source["prompt"];
	        this.rating = source["rating"];
	        this.notes = source["notes"];
	    }
	}
	export class PromptVersion {
	    id: string;
	    savedAt: string;
	    label: string;
	    prompt: string;
	
	    static createFrom(source: any = {}) {
	        return new PromptVersion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.savedAt = source["savedAt"];
	        this.label = source["label"];
	        this.prompt = source["prompt"];
	    }
	}
	export class Variant {
	    id: string;
	    label: string;
	    prompt: string;
	
	    static createFrom(source: any = {}) {
	        return new Variant(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.label = source["label"];
	        this.prompt = source["prompt"];
	    }
	}
	export class ShotSpec {
	    durationSec?: number;
	    fps?: number;
	    aspectRatio?: string;
	    lens?: string;
	    movement?: string;
	    size?: string;
	    angle?: string;
	
	    static createFrom(source: any = {}) {
	        return new ShotSpec(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.durationSec = source["durationSec"];
	        this.fps = source["fps"];
	        this.aspectRatio = source["aspectRatio"];
	        this.lens = source["lens"];
	        this.movement = source["movement"];
	        this.size = source["size"];
	        this.angle = source["angle"];
	    }
	}
	export class Shot {
	    id: string;
	    name: string;
	    intent: string;
	    spec: ShotSpec;
	    prompt: string;
	    lockedLines: number[];
	    mutedLines: number[];
	    beatSheet: BeatDirection[];
	    targetModel?: string;
	    maxChars?: number;
	    variants: Variant[];
	    history: PromptVersion[];
	    takes: Take[];
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Shot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.intent = source["intent"];
	        this.spec = this.convertValues(source["spec"], ShotSpec);
	        this.prompt = source["prompt"];
	        this.lockedLines = source["lockedLines"];
	        this.mutedLines = source["mutedLines"];
	        this.beatSheet = this.convertValues(source["beatSheet"], BeatDirection);
	        this.targetModel = source["targetModel"];
	        this.maxChars = source["maxChars"];
	        this.variants = this.convertValues(source["variants"], Variant);
	        this.history = this.convertValues(source["history"], PromptVersion);
	        this.takes = this.convertValues(source["takes"], Take);
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
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
	export class Scene {
	    id: string;
	    name: string;
	    synopsis: string;
	    shots: Shot[];
	
	    static createFrom(source: any = {}) {
	        return new Scene(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.synopsis = source["synopsis"];
	        this.shots = this.convertValues(source["shots"], Shot);
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
	export class ProjectDefaults {
	    aspectRatio: string;
	    fps: number;
	    durationSec: number;
	    targetModel: string;
	    brain: string;
	    localEndpoint?: string;
	    localModel?: string;
	
	    static createFrom(source: any = {}) {
	        return new ProjectDefaults(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.aspectRatio = source["aspectRatio"];
	        this.fps = source["fps"];
	        this.durationSec = source["durationSec"];
	        this.targetModel = source["targetModel"];
	        this.brain = source["brain"];
	        this.localEndpoint = source["localEndpoint"];
	        this.localModel = source["localModel"];
	    }
	}
	export class Project {
	    id: string;
	    name: string;
	    logline: string;
	    world: string;
	    defaults: ProjectDefaults;
	    scenes: Scene[];
	    characters: CharacterSheet[];
	    artDept: ArtDeptSheet[];
	    locations: LocationSheet[];
	    lookbook: StyleProfile[];
	    references: Reference[];
	    mySetups: CustomSetup[];
	    music?: MusicCue[];
	    voices?: VoiceSheet[];
	    copilot?: ChatMsg[];
	    createdAt: string;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.logline = source["logline"];
	        this.world = source["world"];
	        this.defaults = this.convertValues(source["defaults"], ProjectDefaults);
	        this.scenes = this.convertValues(source["scenes"], Scene);
	        this.characters = this.convertValues(source["characters"], CharacterSheet);
	        this.artDept = this.convertValues(source["artDept"], ArtDeptSheet);
	        this.locations = this.convertValues(source["locations"], LocationSheet);
	        this.lookbook = this.convertValues(source["lookbook"], StyleProfile);
	        this.references = this.convertValues(source["references"], Reference);
	        this.mySetups = this.convertValues(source["mySetups"], CustomSetup);
	        this.music = this.convertValues(source["music"], MusicCue);
	        this.voices = this.convertValues(source["voices"], VoiceSheet);
	        this.copilot = this.convertValues(source["copilot"], ChatMsg);
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
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
	
	export class ProjectMeta {
	    id: string;
	    name: string;
	    logline: string;
	    path: string;
	    updatedAt: string;
	    sceneCount: number;
	    shotCount: number;
	
	    static createFrom(source: any = {}) {
	        return new ProjectMeta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.logline = source["logline"];
	        this.path = source["path"];
	        this.updatedAt = source["updatedAt"];
	        this.sceneCount = source["sceneCount"];
	        this.shotCount = source["shotCount"];
	    }
	}
	
	
	
	
	
	
	
	

}

