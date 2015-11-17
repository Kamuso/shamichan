const admin = require('../server/admin'),
	common = require('../common'),
	util = require('./util')

/**
 * Reads thread and post data from the database
 */
export default class Reader {
	/**
	 * Constructs new database reader
	 * @param {string} board
	 * @param {Object} ident
	 */
	constructor(board, ident) {
		this.ident = ident
		this.board = board
		if (common.checkAuth('janitor', ident)) {
			this.hasAuth = true
			this.canModerate = common.checkAuth('moderator', ident)
		}
	}

	/**
	 * Retrieve thread JSON from the database
	 * @param {int} id - Thread id
	 * @param {Object} opts - Extra options
	 * @returns {(Object|null)} - Retrieved post or null
	 */
	async getThread(id, opts) {

	}

	/**
	 * Read a single post from the database
	 * @param {int} id - Post id
	 * @returns {(Object|null)} - Retrieved post or null
	 */
	async getPost(id) {
		return this.parsePost(await util.getPost(id).run(rcon))
	}

	/**
	 * Adjust post according to the reading client's access priveledges
	 * @param {Object} post - Post object
	 * @returns {(Object|null)} - Parsed post object or null, if client not
	 * 	allowed to view post
	 */
	parsePost(post) {
		if (!post)
			return null
		if (!this.hasAuth) {
			if (post.deleted)
				return null
			if (post.imgDeleted)
				delete post.image
		}
		if (this.canModerate) {
			const mnemonic = admin.genMnemonic(post.ip)
			if (mnemonic)
				post.mnemnic = mnemonic
		}
		util.formatPost(post)
		return post
	}
}
