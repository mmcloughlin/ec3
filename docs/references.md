# References
## Implementation
* [**Selecting Elliptic Curves for Cryptography: An Efficiency and Security Analysis**](https://eprint.iacr.org/2014/130) Joppe W. Bos and Craig Costello and Patrick Longa and Michael Naehrig. _Note:_ Supurbly detailed paper on the implementation of the MSR Elliptic Curve Cryptography Library, with code provided.
* [MSR Elliptic Curve Cryptography Library](https://www.microsoft.com/en-us/research/project/msr-elliptic-curve-cryptography-library/)
* [Elliptic curves and their implementation](https://www.imperialviolet.org/2010/12/04/ecc.html) Adam Langley. _Note:_ Friendly introduction.
* [Multiprecision Arithmetic for Cryptology in C++](https://arxiv.org/abs/1804.07236) ([slides](https://martindale.info/eipsi/ctbignum_slides_2018.pdf), [code](https://github.com/niekbouman/ctbignum)) Niek J. Bouman.
* [The Design and Application of a Retargetable Peephole Optimizer](http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.64.7226&rep=rep1&type=pdf) Davidson and Fraser. _Note:_ Early description of an RTL language.
* [New Instructions Supporting Large Integer Arithmetic on Intel Architecture Processors](https://www.intel.com/content/dam/www/public/us/en/documents/white-papers/ia-large-integer-arithmetic-paper.pdf) Intel White Paper. _Note:_ Concise and easy description of how to exploit the `MULX/ADCX/ADOX` instructions, illustrated with a 512x512-bit multiply.
* [Large Integer Squaring on Intel Architecture Processors](https://www.intel.com/content/dam/www/public/us/en/documents/white-papers/large-integer-squaring-ia-paper.pdf) Intel White Paper. _Note:_ Extension of the multiplcation paper to the squaring case.
* [**Software Implementation of Public-Key Cryptography (SAC Summer School)**](https://irp-cdn.multiscreensite.com/7fa75f95/files/uploaded/Software%20Implementation%20of%20Public-Key%20Cryptography%20-%20SAC%20Summer%20School.pdf) Patrick Longa. _Note:_ Incredible survey.
## Addition Chains
* [Pippenger's exponentiation algorithm](https://cr.yp.to/papers/pippenger.pdf) DJB. _Note:_ Excellent summary. Matrix representation of chains.
* [On addition chains](https://projecteuclid.org/euclid.bams/1183502136) Brauer. _Note:_ The _original_ Brauer paper.
* [The Additive Complexity of a Natural Number](https://www.researchgate.net/publication/267016251_The_additive_complexity_of_a_natural_number) Belaga. _Note:_ Early introduction of the dictionary method, with the observation that only odd window values are required (page 7). Used to prove upper bound.
* [On the Evaluation of Powers](https://www.ii.uni.wroc.pl/~aje/WordEq2015/papers/addition_chains_Yao.pdf) Yao. _Note:_ _Original_ definition of Yao's algorithm.
* [Exponentiating faster with addition chains](https://www.researchgate.net/publication/221348287_Exponentiating_Faster_with_Addition_Chains) Yacobi. _Note:_ LZ applied to dictionary creation. The _original_ definition of Yacobi's algorithm.
* [Efficient computation of addition-subtraction chains using generalized continued Fractions](https://eprint.iacr.org/2013/466) Tall, Sanghare. _Note:_ Adapts continued fractions to subtraction chains. Nice clear review of continued fractions strategies.
* [Addition chains using continued fractions](https://drive.google.com/file/d/1FiCIhWXQl5FY6pcpriAeHQJaBCYZVcvi/view) BBBD. _Note:_ _Original_ continued fractions paper.
* [Calculating optimal addition chains](https://link.springer.com/content/pdf/10.1007%2Fs00607-010-0118-8.pdf)
* [Efficient Generation of Minimal Length Addition Chains](https://pdfs.semanticscholar.org/6e33/657f2acf01c70fb66fbcc9c06416123c7ed6.pdf)
* [Finding short and implementation-friendly addition chains with evolutionary algorithms](https://dspace.mit.edu/bitstream/handle/1721.1/115968/10732_2017_9340_ReferencePDF.pdf)
* [Modifications of Bos and Coster’s Heuristics in search of a shorter addition chain for faster exponentiation](http://library.isical.ac.in:8080/jspui/bitstream/123456789/6441/1/DISS-285.pdf)
* [Developing an automatic generation tool for cryptographic pairing functions](http://doras.dcu.ie/16002/1/thesis.pdf)
* [Addition Chain Heuristics](https://link.springer.com/content/pdf/10.1007/0-387-34805-0_37.pdf)
* [Optimizing linear maps modulo 2](https://binary.cr.yp.to/linearmod2-20090830.pdf)
* [An Artificial Immune System Heuristic for Generating Short Addition Chains](http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.490.991&rep=rep1&type=pdf)
* [SPEED 2009 Proceedings: Software Performance Enhancement for Encryption and Decryption and Cryptographic Compilers](http://www.hyperelliptic.org/SPEED/record09.pdf)
* [Efficient exponentiation using precomputation and vector addition chains](https://link.springer.com/content/pdf/10.1007/BFb0053453.pdf)
* [Implementation of the GBD cryptosystem](https://eprints.qut.edu.au/441/1/08_brown04implementation.pdf)
* [Brauer Addition-Subtraction Chains](http://www.martin-otto.de/publications/docs/2001_MartinOtto_Diplom_BrauerAddition-SubtractionChains.pdf)
* [Speeding up Subgroup Cryptosystems](https://cr.yp.to/bib/2003/stam-thesis.pdf)
* [MpNT: A Multi-Precision Number Theory Package. Number-Theoretic Algorithms.](https://profs.info.uaic.ro/~tr/tr03-02.pdf)
* [Handbook of  Elliptic and  Hyperelliptic  Curve Cryptography: Section 9.2: Fixed Exponent](https://drive.google.com/file/d/10VG_UoYfvhkUon7_I2dGXI9BwQ-MKtND/view)
* [Redundant integer representations and fast exponentiation](https://repository.royalholloway.ac.uk/file/cb0f3ec9-a23f-6ab8-5dd3-1b73129cef71/8/rirafe.pdf)
* [High-speed high-security signatures](https://ed25519.cr.yp.to/ed25519-20110705.pdf)
* [Cultivating Sapling: New Crypto Foundations](https://github.com/zcash/zcash-blog/blob/4268b1a9f7e523f59bb254b9fa7a9f9f9d7f75a7/_posts/2017-07-26-cultivating-sapling-new-crypto-foundations.md) _Note:_ Section "New multi-exponentiation algorithm"
* [Faster batch forgery identification](https://eprint.iacr.org/2012/549) _Note:_ Special case of Pippenger's algorithm, Page 15.
* [On Fast Calculation of Addition Chains for Isogeny-Based Cryptography](https://eprint.iacr.org/2016/1045)
* [Fast and simple constant-time hashing to the BLS12-381 elliptic curve](https://eprint.iacr.org/2019/403)
* [Efficient computation of addition chains](http://www.numdam.org/article/JTNB_1994__6_1_21_0.pdf)
* [New Methods for Generating Short Addition Chains](https://pdfs.semanticscholar.org/b398/d10faca35af9ce5a6026458b251fd0a5640c.pdf)
* [On the final exponentiation for calculating pairings on ordinary elliptic curves](https://eprint.iacr.org/2008/490) _Note:_ Contains various real-world examples of addition sequences.
* [`kwantam/addchain`](https://github.com/kwantam/addchain) _Note:_ Go library for addition chain generation.
* [Addition Chains as Polymorphic Higher-order Functions](https://briansmith.org/addition-chains-as-higher-order-functions-01)
* [**The Most Efficient Known Addition Chains for Field Element & Scalar Inversion for the Most Popular & Most Unpopular Elliptic Curves**](https://briansmith.org/ecc-inversion-addition-chains-01)
* [additionchains.com](http://www.additionchains.com/)
* [OEIS A003313: Length of shortest addition chain for n](https://oeis.org/A003313)
## Finite Field Arithmetic
* [Faster Square Roots in Annoying Finite Fields](https://cr.yp.to/papers/sqroot.pdf) DJB. _Note:_ Methods for square roots modulo primes 1 (mod 4), with focus on the NIST P-224 prime.
* [New Speed Records for Montgomery Modular Multiplication on 8-bit AVR Microcontrollers](https://eprint.iacr.org/2013/882)
* [High-Speed Cryptography](http://www.just.edu.jo/~tawalbeh/cpe776/notes/high.pdf) Cetin Koc. _Note:_ Lecture notes with survey of multiplication methods.
* [High-Speed Algorithms & Architectures For Number-Theoretic Cryptosystems](https://www.microsoft.com/en-us/research/wp-content/uploads/1998/06/97Acar.pdf) Tolga Acar. _Note:_ Thesis. Various proposals for Montgomery multiplies, discussion of how to extend the x86 instruction set for crypto.
* [Generalised Mersenne Numbers Revisited](https://eprint.iacr.org/2011/444)
## Scalar Multiplication/Recoding
* [EE 371 Lecture 5: More Adders & Multipliers](http://web.stanford.edu/class/archive/ee/ee371/ee371.1066/lectures/lect_05.pdf) Mark Horowitz. _Note:_ Nice coverage of Modified Booth Recoding.
* [Scalar-multiplication algorithms](https://cryptojedi.org/peter/data/eccss-20130911b.pdf) Peter Schwabe. _Note:_ ECC 2013 Summer School. Simple introduction to scalar multiplication algorithms.
* [Survey of Elliptic Curve Scalar Multiplication Algorithms](https://pdfs.semanticscholar.org/4da9/cfe2ff561a29bd90c8799873154b819f5cd9.pdf) _Note:_ Concise survey.
* [Fast and Regular Algorithms for Scalar Multiplication over Elliptic Curves](https://eprint.iacr.org/2011/338) Matthieu Rivain. _Note:_ Only skimmed, but the sub-section on "Scalar Recoding" includes a dense summary of the approaches.
* [The Width-w NAF Method Provides Small Memory and Fast Elliptic Scalar Multiplications Secure against Side Channel Attacks](https://drive.google.com/file/d/1VUMq5ef9f4L3kHtKLDPbY4wkIJobLtvz/view?usp=sharing) _Note:_ Signed windows using all _odd_ digits.
* [Securing Elliptic Curve Point Multiplication against Side-Channel Attacks](https://www.bmoeller.de/pdf/ecc-sca-isc2001.pdf) Bodo Moller. _Note:_ Booth-like recoding that _avoids zero_ digits.
## Elliptic Curve Arithmetic
* [Analysis and optimization of elliptic-curve single-scalar multiplication](http://www.hyperelliptic.org/EFD/precomp.pdf) DJB & Lange. _Note:_ the paper behind the Explicit Formula Database. Excellent overview.
* [Complete addition formulas for prime order elliptic curves](https://eprint.iacr.org/2015/1060) Joost Renes and Craig Costello and Lejla Batina.
* [A Compact and Exception-Free Ladder for All Short Weierstrass Elliptic Curves](https://drive.google.com/file/d/1f6K2vvn067hxkTxsAf7Z7W2UNJAbsSB4/view?usp=sharing) Ruggero Susella and Sofia Montrasio.
* [A Binary Redundant Scalar Point Multiplication in Secure Elliptic Curve Cryptosystems](http://ijns.jalaxy.com.tw/contents/ijns-v3-n2/ijns-2006-v3-n2-p132-137.pdf) Sangook Moon. _Note:_ Suggests the use of a _quadruple_ formula in scalar multiplication.
* [Trading Inversions for Multiplications in Elliptic Curve Cryptography](http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.92.1336&rep=rep1&type=pdf) _Note:_ Explores alternative bases and formulae for scalar multiplication.
* [Fast Elliptic Curve Arithmetic and Improved Weil Pairing Evaluation](https://arxiv.org/abs/math/0208038) Kirsten Eisentraeger, Kristin Lauter, Peter L. Montgomery. _Note:_ Merged double-add.
* [Improved Algorithms for Elliptic Curve Arithmetic in GF(2^n)](https://link.springer.com/content/pdf/10.1007/3-540-48892-8_16.pdf) Julio Lopez and Ricardo Dahab. _Note:_ Repeated doubling formulae.
* ["Nice" Curves](https://eprint.iacr.org/2019/1259) Kaushik Nath and Palash Sarkar. _Note:_ Proposes new pairs of Montgomery-Edwards curves at both the 128-bit and the 224-bit security levels. _Note:_ performance boost at the cost of a couple of bits of security. Additional free bits in the top limb allow for savings in modular reduction.
## Bugs
* [Assessing and Exploiting BigNum Vulnerabilities](https://comsecuris.com/slides/slides-bignum-bhus2015.pdf) Ralf-Philipp Weinmann. _Note:_ Examples of specific arbitrary precision integer bugs and how to find them.
* [Practical realisation and elimination of an ECC-related software bug attack](https://eprint.iacr.org/2011/633) B.B. Brumley and M. Barbosa and D. Page and F. Vercauteren. _Note:_ Exploit for bug in OpenSSL version 0.9.8g which permit an attack against ECDH-based functionality. Some discussion of discovery and verification with `CAOverif`.
## Software
* [`zkcrypto/jubjub`](https://github.com/zkcrypto/jubjub)
## Code Generation
* [Fast and simple constant-time hashing to the BLS12-381 elliptic curve](https://eprint.iacr.org/2019/403)
* [Designing a code generator of Pairing Based Cryptography functions](https://www.ucc.ie/en/media/academic/centreforplanningeducationresearch/LDominguez.pdf)
* [Ideas for a New Elliptic Curve Cryptography Library](https://briansmith.org/GFp-0)
## Uncategorized
* [A Brief Overview of HOL4](https://ts.data61.csiro.au/publications/nicta_full_text/1482.pdf)
* [On Construction of a Library of Formally Verified Low-level Arithmetic Functions](https://staff.aist.go.jp/reynald.affeldt/documents/arilib-affeldt.pdf)
* [Verification of Machine Code Implementations of Arithmetic Functions for Cryptography](http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.171.621&rep=rep1&type=pdf)
* [A Verified Extensible Library of Elliptic Curves](https://hal.inria.fr/hal-01425957/document) Jean Karim Zinzindohoue, Evmorfia-Iro Bartzia, Karthikeyan Bhargavan.
* [Formal Verification of a Big Integer Library Written in C0](http://www-wjp.cs.uni-sb.de/publikationen/Fi06.pdf) Sabine Fischer.
* [Proof Pearl: A Verified Bignum Implementation in x86-64 Machine Code](https://www.cl.cam.ac.uk/~mom22/cpp13/) Magnus O. Myreen and Gregorio Curello.
* [Reconstruction of Z3's Bit-Vector Proofs in HOL4 and Isabelle/HOL](http://user.it.uu.se/~tjawe125/publications/boehme11reconstruction.html)
* [A Why3 Framework for Reflection Proofs and its Application to GMP's Algorithms](https://hal.inria.fr/hal-01699754v1/document)
* [How to Get an Efficient yet Verified Arbitrary-Precision Integer Library](https://hal.inria.fr/hal-01519732v2/document)
* [Analyzing and Comparing Montgomery Multiplication Algorithms](https://pdfs.semanticscholar.org/5e39/41ff482ec3ee41dc53c3298f0be085c69483.pdf) Cetin K. Koc, Tolga Acar, Burton S. Kaliski Jr.
* [Optimizing Multiprecision Multiplication for Public Key Cryptography](https://eprint.iacr.org/2007/299) Michael Scott and Piotr Szczechowiak.
* [Fast Multi-Precision Multiplication for Public-Key Cryptography on Embedded Microprocessors](https://www.iacr.org/archive/ches2011/69170459/69170459.pdf) Michael Hutter and Erich Wenger.
* [A Verified, Efficient Embedding of a Verifiable Assembly Language](https://www.microsoft.com/en-us/research/publication/a-verified-efficient-embedding-of-a-verifiable-assembly-language/)
* [Twisted Edwards Curves Revisited](https://eprint.iacr.org/2008/522) Huseyin Hisil and Kenneth Koon-Ho Wong and Gary Carter and Ed Dawson.
* [CAOVerif: An Open-Source Deductive Verification Platform for Cryptographic Software Implementations](https://haslab.uminho.pt/jsp/files/opencertjournal_ack.pdf)
* [Efficient Arithmetic In (Pseudo-)Mersenne Prime Order Fields](https://eprint.iacr.org/2018/985) Kaushik Nath and Palash Sarkar.
* [A type-safe arbitrary precision arithmetic portability layer for HLS tools](https://hal.inria.fr/hal-02131798v2/document)
* [Systematic Support of Parallel Bit Streams in LLVM](https://pdfs.semanticscholar.org/1530/6ba903de3e674ba6facf9a76ea41880fc918.pdf) Meng Lin.
* [Montgomery Multiplication Using Vector Instructions](https://eprint.iacr.org/2013/519) ([slides](http://sac2013.irmacs.sfu.ca/slides/s26.pdf)) Joppe W. Bos and Peter L. Montgomery and Daniel Shumow and Gregory M. Zaverucha.
* [Complete Systems of Two Addition Laws for Elliptic Curves](https://drive.google.com/file/d/12RLNZz1MdDXuBL4wTU8V33Ff6dFCoVRx/view?usp=sharing) Bosma W. Lenstra H.W.
* [Crafting Certified Elliptic Curve Cryptography Implementations in Coq](http://adam.chlipala.net/theses/andreser_meng.pdf)
* [Specification of Curve Selection and Supported Curve Parameters in MSR ECCLib](https://www.microsoft.com/en-us/research/publication/specification-of-curve-selection-and-supported-curve-parameters-in-msr-ecclib/)
* [Scalar recoding and regular 2w-ary right-to-left EC scalar multiplication algorithm](https://drive.google.com/open?id=1oRbgAyqN8H3ne8nDBhbvCA4NkyE9zi6z) Yoo-Jin Baek.
* [Signed Binary Representations Revisited](https://eprint.iacr.org/2004/195) Katsuyuki Okeya and Katja Schmidt-Samoa and Christian Spahn and Tsuyoshi Takagi.
* [Some Explicit Formulae of NAF and its Left-to-Right Analogue](https://eprint.iacr.org/2005/384) Dong-Guk Han and Tetsuya Izu and Tsuyoshi Takagi.
* [Software Implementation of Modular Exponentiation Using Advanced Vector Instructions Architectures](https://drive.google.com/file/d/1l9-ZEOl9-cgVUjTwkTeI7YyrEid1y4BG/view?usp=sharing) Shay Gueron and Vlad Krasnov.
* [Simple High-Level Code For Cryptographic Arithmetic – With Proofs, Without Compromises](http://adam.chlipala.net/papers/FiatCryptoSP19/FiatCryptoSP19.pdf) _Note:_ FIAT Crypto
* [Exponent Recoding and Regular Exponentiation Algorithms](http://www.geocities.ws/mike.tunstall/papers/JT09.pdf) Joye-Tunstall.
* [Efficient arithmetic on Koblitz curves](https://www.decred.org/research/solinas2000.pdf) Solinas.
* [A survey of fast exponentiation methods](https://www.dmgordon.org/papers/jalg.pdf) Daniel M. Gordon.
* [Parallel scalar multiplication on general elliptic curves over F_p hedged against Non-Differential Side-Channel Attacks](https://pdfs.semanticscholar.org/ffa2/f1db6aeaf0fb03cf010bad12266da11d00e1.pdf) Wieland Fischer, Christophe Giraud, Erik Woodward Knudsen and Jean-Pierre Seifert.
* [Speeding the Pollard and Elliptic Curve Methods of Factorization](https://www.ams.org/journals/mcom/1987-48-177/S0025-5718-1987-0866113-7/S0025-5718-1987-0866113-7.pdf) Peter Montgomery.
* [Highly Regular Right-to-Left Algorithms for Scalar Multiplication](https://www.iacr.org/archive/ches2007/47270135/47270135.pdf) Marc Joye.
* [Slothful reduction](https://eprint.iacr.org/2017/437) Michael Scott.
* [Faster addition and doubling on elliptic curves](https://eprint.iacr.org/2007/286) Daniel J. Bernstein and Tanja Lange. _Note:_ Found via citation "Thankfully efficient exception-free formulations are now available [4], [11]".
* [RELIC Toolkit](https://github.com/relic-toolkit/relic)
* [Efficient Elliptic Curve Exponentiation Using Mixed Coordinates](https://link.springer.com/content/pdf/10.1007/3-540-49649-1_6.pdf)
* [Improved Techniques for Fast Exponentiation](https://www.bmoeller.de/pdf/fastexp-icisc2002.pdf)
* [Finite Field Arithmetic (Chapter 11 of Handbook of Elliptic and Hyperelliptic Curve Cryptogrpahy)](https://koclab.cs.ucsb.edu/teaching/ecc/eccPapers/Doche-ch11.pdf) Doche.
* [Go assembly language complementary reference](https://quasilyte.dev/blog/post/go-asm-complementary-reference/) Iskander Sharipov.
* [Fast point multiplication algorithms for binary elliptic curves with and without precomputation](https://eprint.iacr.org/2014/427)
* [Addition chains with multiplicative cost](http://www.math.ucsd.edu/~ronspubs/78_11_addition_chains.pdf) _Note:_ Theoretical result on number of multiplications.
* [On the evaluation of powers and related problems](https://drive.google.com/file/d/1vi2nO59glqWSgINX85yerWwW1omrktMY/view?usp=sharing) Pippenger. _Note:_ Matrix and graph-theoretical view of addition chain methods.
* [On string replacement exponentiation](https://drive.google.com/file/d/1zpbb0ZJzewMJGqKbT_pRcS7MXscDcglJ/view) _Note:_ Mostly just theoretical results.
* [High-radix and bit recoding techniques for modular exponentiation](https://drive.google.com/file/d/1TD1v7E9qoFW5J2vt8bntmdQBLMw-VItU/view) _Note:_ Requires subtraction.
* [Fast Modular Exponentiation](http://cryptocode.net/docs/c06.pdf) _Note:_ Requires subtraction.
* [An improved binary algorithm for RSA](https://core.ac.uk/download/pdf/82735567.pdf) _Note:_ Requires subtraction.
* [New Algorithm for Classical Modular Inverse](https://link.springer.com/content/pdf/10.1007%2F3-540-36400-5_6.pdf) _Note:_ Optimized for hardware, requires subtraction.
* [Algorithms for software implementations of RSA](http://www.chrismitchell.net/afsior.pdf) _Note:_ Lookup-table based.
* [SafeCurves: choosing safe curves for elliptic-curve cryptography](https://safecurves.cr.yp.to/)
* [Fast Prime Field Elliptic Curve Cryptography with 256 Bit Primes](https://eprint.iacr.org/2013/816)
* [Modular Multiplication Without Trial Division](https://web.itu.edu.tr/~orssi/dersler/cryptography/Montgomery.pdf)
* [The Montgomery Modular Inverse - Revisited](https://drive.google.com/file/d/1VGTCc9EasatNYM6rS2xkVahX1YFzkxX6/view)
* [Efficient Software-Implementation of Finite Fields with Applications to Cryptography](https://drive.google.com/file/d/1sUJyBEtLTGVnqotHSZQBgbTecWtnDU1w/view)
* [Efficient Software Implementations of Modular Exponentiation](https://eprint.iacr.org/2011/239)
* [The Montgomery Powering Ladder](http://cr.yp.to/bib/2003/joye-ladder.pdf)
* [A note on how to (pre-)compute a ladder](https://eprint.iacr.org/2017/264)
* [A Faster Software Implementation of the Supersingular Isogeny Diffie-Hellman Key Exchange Protocol](https://eprint.iacr.org/2017/1015)
* [Fast Arithmetic Modulo 2^xp^y±1](https://eprint.iacr.org/2016/986)
* [RFC 7748: Elliptic Curves for Security](https://www.rfc-editor.org/rfc/rfc7748.txt)
* [Generalized Mersenne Numbers](http://cacr.uwaterloo.ca/techreports/1999/corr99-39.pdf)
* [Solinas primes of small weight for fixed sizes](https://eprint.iacr.org/2010/058)
* [Twisted Edwards Curves](https://eprint.iacr.org/2008/013)
* [New software speed records for cryptographic pairings](https://cryptojedi.org/papers/dclxvi-20100714.pdf)
* [Multiplication and Squaring on Pairing-Friendly Fields](https://eprint.iacr.org/2006/471)
* [Efficient Arithmetic Modulo Minimal Redundancy Cyclotomic Primes](https://pdfs.semanticscholar.org/45eb/cc684c6c5b1d0941a490280382012f1b2361.pdf)
* [Pairing-Friendly Curves](https://tools.ietf.org/id/draft-yonezawa-pairing-friendly-curves-00.html)
* [Aspects of Hyperelliptic Curves over Large Prime Fields in Software Implementations](https://www.iacr.org/archive/ches2004/31560147/31560147.pdf)
* [Software implementation of the NIST elliptic curves over prime fields](http://delta.cs.cinvestav.mx/~francisco/arith/julio.pdf)
* [Speeding up XTR](https://www.iacr.org/archive/asiacrypt2001/22480125.pdf)
* [Engineering Primes: Taking the Magic out of Magic Numbers](https://speakerdeck.com/gtank/engineering-prime-numbers)
* [How to (pre-)compute a ladder](http://www.ic.unicamp.br/~ra142685/userfiles/papers/oliveira_sac2017.pdf)
* [Double-Base Number System and Applications](http://www.hyperelliptic.org/tanja/conf/ECC08/slides/Christophe-Doche.pdf)
* [An Introduction & Supplement to Knuth's Introduction to Addition Chains](https://briansmith.org/addition-chain-intro-01)
* [Curve25519: new Diffie-Hellman speed records](https://cr.yp.to/ecdh/curve25519-20060209.pdf)
* [Ed25519: High-speed high-security signatures](https://ed25519.cr.yp.to/ed25519-20110926.pdf)
* [Ed448-Goldilocks, a new elliptic curve](https://eprint.iacr.org/2015/625)
* [FourQ: four-dimensional decompositions on a Q-curve over the Mersenne prime](https://eprint.iacr.org/2015/565)
* [Curve4Q: `draft-ladd-cfrg-4q-01`](https://tools.ietf.org/html/draft-ladd-cfrg-4q-01)
* [Jacobian coordinates with a4=-3 for short Weierstrass curves](https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html)
* [FIPS 186-2](https://csrc.nist.gov/csrc/media/publications/fips/186/2/archive/2000-01-27/documents/fips186-2.pdf)
* [SEC2](https://safecurves.cr.yp.to/www.secg.org/sec2-v2.pdf)
* [SIKE/SIDH Spec](https://sike.org/files/SIDH-spec.pdf)
* [Handbook of Applied Cryptography: Chapter 14: Efficient Implementation](http://cacr.uwaterloo.ca/hac/about/chap14.pdf)
* [Handbook of Elliptic and Hyperelliptic Curve Cryptography](https://drive.google.com/file/d/10VG_UoYfvhkUon7_I2dGXI9BwQ-MKtND/view)
* [Guide to Elliptic Curve Cryptography](https://drive.google.com/file/d/1_2t06CvYwTbZcvrn6nZ_ms3O__0xRl9B/view)
* [Topics in Computational Number Theory Inspired by Peter L. Montgomery](http://www.joppebos.com/_pages/montgomery/)
