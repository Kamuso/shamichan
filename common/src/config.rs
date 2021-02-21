use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Upload size constraints
#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct UploadMaximums {
	/// Max size in MB
	pub size: f64,

	/// Max width in pixels
	pub width: u64,

	/// Max height in pixels
	pub height: u64,
}

impl Default for UploadMaximums {
	#[inline]
	fn default() -> Self {
		Self {
			size: 5.0,
			width: 600,
			height: 600,
		}
	}
}

/// Upload configurations
#[derive(Serialize, Deserialize, Default, Debug, Clone)]
pub struct Uploads {
	/// Use JPEG thumbnails instead of WEBP
	pub jpeg_thumbnails: bool,

	/// Upload size constraints
	pub max: UploadMaximums,
}

/// Available user interface languages
#[allow(non_camel_case_types)]
#[derive(Serialize, Deserialize, Hash, Eq, PartialEq, Debug, Clone)]
pub enum Language {
	en_GB,
	es_ES,
	fr_FR,
	nl_NL,
	pl_PL,
	pt_BR,
	ru_RU,
	sk_SK,
	tr_TR,
	uk_UA,
	zh_TW,
}

impl Default for Language {
	#[inline]
	fn default() -> Self {
		Self::en_GB
	}
}

impl std::fmt::Display for Language {
	#[inline]
	fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
		f.write_str(match self {
			Self::en_GB => "en_GB",
			Self::es_ES => "es_ES",
			Self::fr_FR => "fr_FR",
			Self::nl_NL => "nl_NL",
			Self::pl_PL => "pl_PL",
			Self::pt_BR => "pt_BR",
			Self::ru_RU => "ru_RU",
			Self::sk_SK => "sk_SK",
			Self::tr_TR => "tr_TR",
			Self::uk_UA => "uk_UA",
			Self::zh_TW => "zh_TW",
		})
	}
}

/// Available user interface themes
#[allow(non_camel_case_types)]
#[derive(Serialize, Deserialize, Debug, Clone)]
pub enum Theme {
	ashita,
	console,
	egophobe,
	gar,
	glass,
	gowno,
	higan,
	inumi,
	mawaru,
	moe,
	moon,
	ocean,
	rave,
	tavern,
	tea,
	win95,
}

impl Default for Theme {
	#[inline]
	fn default() -> Self {
		Self::ashita
	}
}

impl std::fmt::Display for Theme {
	fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
		f.write_str(match self {
			Self::ashita => "ashita",
			Self::console => "console",
			Self::egophobe => "egophobe",
			Self::gar => "gar",
			Self::glass => "glass",
			Self::gowno => "gowno",
			Self::higan => "higan",
			Self::inumi => "inumi",
			Self::mawaru => "mawaru",
			Self::moe => "moe",
			Self::moon => "moon",
			Self::ocean => "ocean",
			Self::rave => "rave",
			Self::tavern => "tavern",
			Self::tea => "tea",
			Self::win95 => "win95",
		})
	}
}

/// Global server configurations exposed to the client
#[derive(Serialize, Deserialize, Default, Debug, Clone)]
pub struct Public {
	/// Mark site content for mature audiences
	pub mature: bool,

	/// Enable captchas and antispam
	pub enable_antispam: bool,

	/// Delete unused threads
	pub prune_threads: bool,

	/// Days a thread stays unpruned without bumping.
	/// 0 means threads do not expire.
	pub thread_expiry: u32,

	/// Default client interface language
	pub default_lang: Language,

	/// Default client interface theme
	pub default_theme: Theme,

	/// Configured labeled links to resources
	pub links: HashMap<String, String>,

	/// Info custom information display per language.
	///
	/// If the selected language does not have an entry, the default_lang entry
	/// will be used.
	//
	// TODO: automatically generate the header of info from public configs and
	// language pack template (replace all `\n` with `<br>`)
	pub information: HashMap<Language, String>,

	/// Support email address
	pub support_email: Option<String>,

	/// Upload configurations
	pub uploads: Uploads,
}
