export namespace main {
	
	export class FileReadResult {
	    success: boolean;
	    message: string;
	    bits: string;
	    filePath: string;
	    fileSize: number;
	
	    static createFrom(source: any = {}) {
	        return new FileReadResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.bits = source["bits"];
	        this.filePath = source["filePath"];
	        this.fileSize = source["fileSize"];
	    }
	}
	export class InputValidationResult {
	    valid: boolean;
	    extractedBits: string;
	    bitsCount: number;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new InputValidationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.valid = source["valid"];
	        this.extractedBits = source["extractedBits"];
	        this.bitsCount = source["bitsCount"];
	        this.message = source["message"];
	    }
	}
	export class KeyValidationResult {
	    valid: boolean;
	    binaryKey: string;
	    message: string;
	    keyLength: number;
	
	    static createFrom(source: any = {}) {
	        return new KeyValidationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.valid = source["valid"];
	        this.binaryKey = source["binaryKey"];
	        this.message = source["message"];
	        this.keyLength = source["keyLength"];
	    }
	}
	export class OperationResult {
	    success: boolean;
	    message: string;
	    cipherText: string;
	    keyStream: string;
	    bitsCount: number;
	    extractedBits: string;
	
	    static createFrom(source: any = {}) {
	        return new OperationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.cipherText = source["cipherText"];
	        this.keyStream = source["keyStream"];
	        this.bitsCount = source["bitsCount"];
	        this.extractedBits = source["extractedBits"];
	    }
	}

}

