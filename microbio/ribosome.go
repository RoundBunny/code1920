package microbio

// Ribosome : Here we attempt to model the Translation phase
// of gene expression done by the Ribosome
type Ribosome struct {
}

var startCodon = MakeStrand([]byte("AUG"), RIBOSE)

// Translate returns an AminoAcidStrand for a given NucleotideStrand
// Use WaitForTRNAAndGetAttachedAminoAcid(codon) to get AminoAcid
func (t *Ribosome) Translate(nucleotides NucleotideStrand) AminoAcidStrand {
	
	//declares AA chain
	var acid AminoAcidStrand
	
	//flag variable
	isAdding := false

	//loops through nucleotide strand
	for i:=0; i < nucleotides.Length(); i++ {
		//Gets the values of the current Codon
		cCodon := nucleotides.Codon()

		//If there is no STOP code, return empty 
		if cCodon == nil {
			break;
		}

		//Gets value of current Amino Acid
		cAmino := WaitForTRNAAndGetAttachedAminoAcid(cCodon)

		//If the Amino Acid is a stop code, returns acid chain
		if cAmino == nil {
			return acid
		} 

		//If looking for the start Codon, go one by one until found
		//Once found, sets isAdding to true to being the normal pattern and moves to next codon
		//Normal pattern appends the Amino Acid to the chain and moves onto next codon.
		if isAdding {
			acid = append(acid, cAmino)
			for j:=0; j<3; j++ {nucleotides.SlideRight()}
		} else {
			if cCodon.Matches(startCodon) {
				isAdding = true
				for j:=0; j<3; j++ {nucleotides.SlideRight()}
			} else {
				nucleotides.SlideRight()
			}
		}
		
		
	}

	//returns empty AA chain
	var acidEmpty AminoAcidStrand
	return acidEmpty
}
