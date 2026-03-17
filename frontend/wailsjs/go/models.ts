export namespace main {
	
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
	    binaryKey: string;
	    keyStream: string;
	    originalSize: number;
	    processedSize: number;
	    outputPath: string;
	
	    static createFrom(source: any = {}) {
	        return new OperationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.binaryKey = source["binaryKey"];
	        this.keyStream = source["keyStream"];
	        this.originalSize = source["originalSize"];
	        this.processedSize = source["processedSize"];
	        this.outputPath = source["outputPath"];
	    }
	}

}

