// Package types mirrors the shared TypeScript data model from Slate 0.3.2.
// JSON field names match the original project.json schema for file compatibility.
package types

// ---- Domain model ----

type ShotSpec struct {
	DurationSec *float64 `json:"durationSec"`
	FPS         *float64 `json:"fps"`
	AspectRatio *string  `json:"aspectRatio"`
	Lens        *string  `json:"lens"`
	Movement    *string  `json:"movement"`
	Size        *string  `json:"size"`
	Angle       *string  `json:"angle"`
}

type BeatDirection struct {
	From float64 `json:"from"`
	To   float64 `json:"to"`
	Text string  `json:"text"`
}

type PromptVersion struct {
	ID      string `json:"id"`
	SavedAt string `json:"savedAt"`
	Label   string `json:"label"`
	Prompt  string `json:"prompt"`
}

type Take struct {
	ID      string `json:"id"`
	LoggedAt string `json:"loggedAt"`
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	Rating  string `json:"rating"`
	Notes   string `json:"notes"`
}

type Variant struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Prompt string `json:"prompt"`
}

type Shot struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Intent      string           `json:"intent"`
	Spec        ShotSpec         `json:"spec"`
	Prompt      string           `json:"prompt"`
	LockedLines []int            `json:"lockedLines"`
	MutedLines  []int            `json:"mutedLines"`
	BeatSheet   []BeatDirection  `json:"beatSheet"`
	TargetModel *string          `json:"targetModel"`
	MaxChars    *int             `json:"maxChars"`
	Variants    []Variant        `json:"variants"`
	History     []PromptVersion  `json:"history"`
	Takes       []Take           `json:"takes"`
	CreatedAt   string           `json:"createdAt"`
	UpdatedAt   string           `json:"updatedAt"`
}

type Scene struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Synopsis string `json:"synopsis"`
	Shots    []Shot `json:"shots"`
}

type CharacterSheet struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Age           string   `json:"age"`
	Gender        string   `json:"gender"`
	Ethnicity     string   `json:"ethnicity"`
	FaceFeatures  string   `json:"faceFeatures"`
	Hair          string   `json:"hair"`
	Clothing      string   `json:"clothing"`
	Expression    string   `json:"expression"`
	EyeDirection  string   `json:"eyeDirection"`
	Mood          string   `json:"mood"`
	Environment   string   `json:"environment"`
	KeyLightSide  string   `json:"keyLightSide"`
	LightingMood  string   `json:"lightingMood"`
	Scenario      string   `json:"scenario"`
	Notes         string   `json:"notes"`
	Images        []string `json:"images,omitempty"`
}

type ArtDeptSheet struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Materials   string `json:"materials"`
	Condition   string `json:"condition"`
	Era         string `json:"era"`
	Distinctive string `json:"distinctive"`
	Notes       string `json:"notes"`
}

type LocationSheet struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	InteriorExterior string   `json:"interiorExterior"`
	Description      string   `json:"description"`
	TimeOfDay        string   `json:"timeOfDay"`
	Weather          string   `json:"weather"`
	Architecture     string   `json:"architecture"`
	Textures         string   `json:"textures"`
	PracticalLights  string   `json:"practicalLights"`
	Notes            string   `json:"notes"`
	Images           []string `json:"images,omitempty"`
}

type StyleProfile struct {
	ID           string   `json:"id"`
	Source       string   `json:"source"`
	Kind         string   `json:"kind"`
	Tone         string   `json:"tone"`
	Palette      string   `json:"palette"`
	Lighting     string   `json:"lighting"`
	LensLanguage string   `json:"lensLanguage"`
	Movement     string   `json:"movement"`
	Blocking     string   `json:"blocking"`
	Editorial    string   `json:"editorial"`
	Notes        string   `json:"notes"`
	Images       []string `json:"images,omitempty"`
}

type ElementSheet struct {
	Lensing     string `json:"lensing"`
	Lighting    string `json:"lighting"`
	Palette     string `json:"palette"`
	Composition string `json:"composition"`
	Movement    string `json:"movement"`
	Texture     string `json:"texture"`
	Mood        string `json:"mood"`
	Notes       string `json:"notes"`
}

type CircledTake struct {
	Project   string   `json:"project"`
	Shot      *string  `json:"shot"`
	MediaPath string   `json:"mediaPath"`
	FileName  string   `json:"fileName"`
	Rating    float64  `json:"rating"`
	InSec     *float64 `json:"inSec"`
	OutSec    *float64 `json:"outSec"`
}

type Reference struct {
	ID       string        `json:"id"`
	Path     string        `json:"path"`
	Kind     string        `json:"kind"`
	Label    string        `json:"label"`
	Frames   []string      `json:"frames"`
	Elements *ElementSheet `json:"elements"`
	AddedAt  string        `json:"addedAt"`
}

type CustomSetup struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Snippet  string `json:"snippet"`
	Section  string `json:"section"`
	Tags     []string `json:"tags"`
	Favorite bool   `json:"favorite"`
}

type ProjectDefaults struct {
	AspectRatio   string `json:"aspectRatio"`
	FPS           float64 `json:"fps"`
	DurationSec   float64 `json:"durationSec"`
	TargetModel   string `json:"targetModel"`
	Brain         string `json:"brain"`
	LocalEndpoint string `json:"localEndpoint,omitempty"`
	LocalModel    string `json:"localModel,omitempty"`
}

type AudioFingerprint struct {
	DurationSec           float64  `json:"durationSec"`
	SampledSec            float64  `json:"sampledSec"`
	BpmEstimate           *float64 `json:"bpmEstimate"`
	BpmConfidence         string   `json:"bpmConfidence"`
	PitchMedianHz         *float64 `json:"pitchMedianHz"`
	PitchSpreadSemitones  *float64 `json:"pitchSpreadSemitones"`
	VoicedRatio           float64  `json:"voicedRatio"`
	Brightness            string   `json:"brightness"`
	DynamicRangeDb        float64  `json:"dynamicRangeDb"`
	EnergyArc             string   `json:"energyArc"`
	SilenceRatio          float64  `json:"silenceRatio"`
	LongestSilenceSec     float64  `json:"longestSilenceSec"`
}

type MusicCue struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	SceneRef        string   `json:"sceneRef"`
	Intent          string   `json:"intent"`
	Genre           string   `json:"genre"`
	Mood            string   `json:"mood"`
	Tempo           string   `json:"tempo"`
	Instrumentation string   `json:"instrumentation"`
	Era             string   `json:"era"`
	Structure       string   `json:"structure"`
	Vocals          string   `json:"vocals"`
	LyricTheme      string   `json:"lyricTheme"`
	Lyrics          string   `json:"lyrics"`
	DurationSec     *float64 `json:"durationSec"`
	Notes           string   `json:"notes"`
}

type VoiceSheet struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	CharacterID    *string `json:"characterId"`
	AgeGender      string  `json:"ageGender"`
	Accent         string  `json:"accent"`
	Timbre         string  `json:"timbre"`
	Pitch          string  `json:"pitch"`
	Pacing         string  `json:"pacing"`
	Energy         string  `json:"energy"`
	Texture        string  `json:"texture"`
	EmotionalRange string  `json:"emotionalRange"`
	SampleLine     string  `json:"sampleLine"`
	Notes          string  `json:"notes"`
}

type ChatMsg struct {
	Role     string   `json:"role"`
	Text     string   `json:"text"`
	Receipts []string `json:"receipts,omitempty"`
}

type Project struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Logline    string           `json:"logline"`
	World      string           `json:"world"`
	Defaults   ProjectDefaults  `json:"defaults"`
	Scenes     []Scene          `json:"scenes"`
	Characters []CharacterSheet `json:"characters"`
	ArtDept    []ArtDeptSheet   `json:"artDept"`
	Locations  []LocationSheet  `json:"locations"`
	Lookbook   []StyleProfile   `json:"lookbook"`
	References []Reference      `json:"references"`
	MySetups   []CustomSetup    `json:"mySetups"`
	Music      []MusicCue       `json:"music,omitempty"`
	Voices     []VoiceSheet     `json:"voices,omitempty"`
	Copilot    []ChatMsg        `json:"copilot,omitempty"`
	CreatedAt  string           `json:"createdAt"`
	UpdatedAt  string           `json:"updatedAt"`
}

type ProjectMeta struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Logline    string `json:"logline"`
	Path       string `json:"path"`
	UpdatedAt  string `json:"updatedAt"`
	SceneCount int    `json:"sceneCount"`
	ShotCount  int    `json:"shotCount"`
}

// ---- Brain ----

type BrainRequest struct {
	ID            string   `json:"id"`
	Task          string   `json:"task"`
	System        string   `json:"system"`
	Prompt        string   `json:"prompt"`
	Images        []string `json:"images,omitempty"`
	Tier          string   `json:"tier"`
	ExpectJSON    bool     `json:"expectJson,omitempty"`
	LocalEndpoint string   `json:"localEndpoint,omitempty"`
	LocalModel    string   `json:"localModel,omitempty"`
	Backend       string   `json:"backend,omitempty"`
}

type BrainResult struct {
	ID        string      `json:"id"`
	OK        bool        `json:"ok"`
	Text      string      `json:"text"`
	JSON      interface{} `json:"json,omitempty"`
	Error     string      `json:"error,omitempty"`
	ElapsedMs int64       `json:"elapsedMs"`
}

type BackendStatus struct {
	Available bool   `json:"available"`
	Version   *string `json:"version"`
}

type LocalBackendStatus struct {
	Available bool    `json:"available"`
	Version   *string `json:"version"`
	Endpoint  *string `json:"endpoint"`
}

type BrainStatus struct {
	Claude BackendStatus      `json:"claude"`
	Codex  BackendStatus      `json:"codex"`
	Local  LocalBackendStatus `json:"local"`
}

type LocalModelInfo struct {
	ID string `json:"id"`
}

type LocalModelsResult struct {
	Endpoint *string          `json:"endpoint"`
	Models   []LocalModelInfo `json:"models"`
}

type MediaIngestResult struct {
	Kind   string   `json:"kind"`
	Frames []string `json:"frames"`
}

type LocalOpts struct {
	Endpoint string `json:"endpoint,omitempty"`
	Model    string `json:"model,omitempty"`
}
