use common::payloads::post_body::{Node, PatchNode};

// TODO: keep original text for expensive to parse fragments that can not fail
// (like code tags) and first perform a simple source string comparison to
// ascertain regeneration is in fact needed.

/// Diff the new post body against the old.
///
/// All performed on one thread to maximize thread locality.
/// Yields of work sharing here are doubtable.
pub fn diff(old: &Node, new: &Node) -> Option<PatchNode> {
	use Node::*;

	match (old, new) {
		(Empty, Empty) | (NewLine, NewLine) => None,
		(Siblings(old), Siblings(new)) => {
			macro_rules! diff {
				($i:expr) => {
					diff(&*old[$i], &*new[$i])
				};
			}

			match (diff!(0), diff!(1)) {
				(None, None) => None,
				(l @ _, r @ _) => Some(PatchNode::Siblings([
					l.map(|p| p.into()),
					r.map(|p| p.into()),
				])),
			}
		}
		(Text(new), Text(old))
		| (URL(old), URL(new))
		| (Code(old), Code(new)) => diff_text(old, new),
		(Spoiler(old), Spoiler(new))
		| (Bold(old), Bold(new))
		| (Italic(old), Italic(new))
		| (Quoted(old), Quoted(new)) => diff(old, new),
		(old @ _, new @ _) => {
			if old != new {
				Some(PatchNode::Replace(new.clone()))
			} else {
				None
			}
		}
	}
}

/// Diff text and return patching instructions to enable at least some
/// differential compression for string updates
fn diff_text(old: &str, new: &str) -> Option<PatchNode> {
	// Hot path - most strings won't change and this will compare by length
	// first anyway
	if old == new {
		None
	} else {
		Some(PatchNode::Patch(
			common::payloads::post_body::TextPatch::new(
				&old.chars().collect::<Vec<char>>(),
				&new.chars().collect::<Vec<char>>(),
			),
		))
	}
}
